package db

import (
	"github.com/gin-gonic/gin"
	"github.com/laydong/toolpkg/log"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"strings"
	"time"
)

const (
	defaultPoolMaxIdle     = 2                                 // 连接池空闲连接数量
	defaultPoolMaxOpen     = 13                                // 连接池最大连接数量4c*2+4只读副本+1主实例
	defaultConnMaxLifeTime = time.Second * time.Duration(7200) // MySQL默认长连接时间为8个小时,可根据高并发业务持续时间合理设置该值
	defaultConnMaxIdleTime = time.Second * time.Duration(600)  // 设置连接10分钟没有用到就断开连接(内存要求较高可降低该值)
	LevelInfo              = "info"
	LevelWarn              = "warn"
	LevelError             = "error"
	dbWriteTime            = 3 * time.Second //使用主库时间
)

type DbPoolCfg struct {
	MaxIdleConn int `json:"max_idle_conn"` //空闲连接数
	MaxOpenConn int `json:"max_open_conn"` //最大连接数
	MaxLifeTime int `json:"max_life_time"` //连接可重用的最大时间
	MaxIdleTime int `json:"max_idle_time"` //在关闭连接之前,连接可能处于空闲状态的最大时间
}

var DB *gorm.DB

// InitDB init db
func InitDB(dsn string, dsn1 ...string) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: log.Default(logger.Info)})
	if err != nil {
		//log.Printf("[app.gstore] mysql open fail, err=%s", err)
		panic(err)
	}
	if len(dsn1) > 0 {
		//主从库
		db.Use(dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{mysql.Open(dsn)},
			Replicas: []gorm.Dialector{mysql.Open(dsn1[0])},
			Policy:   dbresolver.RandomPolicy{},
		}).SetConnMaxIdleTime(defaultPoolMaxIdle).
			SetConnMaxLifetime(defaultConnMaxLifeTime).
			SetMaxIdleConns(defaultPoolMaxIdle).
			SetMaxOpenConns(defaultPoolMaxOpen))
	} else {
		//单库
		db.Use(dbresolver.Register(dbresolver.Config{
			Sources: []gorm.Dialector{mysql.Open(dsn)},
			Policy:  dbresolver.RandomPolicy{},
		}).SetConnMaxIdleTime(defaultPoolMaxIdle).
			SetConnMaxLifetime(defaultConnMaxLifeTime).
			SetMaxIdleConns(defaultPoolMaxIdle).
			SetMaxOpenConns(defaultPoolMaxOpen))
	}

	err = Initialize(db)
	//执行sql 主从
	registerReplicaCallbacks(db)
	if err != nil {
		return
	}
	err = DbSurvive(db)
	if err != nil {
		//log.Printf("[app.gstore] mysql survive fail, err=%s", err)
		panic(err)
	}
	DB = db
	//log.Printf("[app.gstore] mysql success")
}

func GetDB(c *gin.Context) *gorm.DB {
	key := viper.GetString("mysql.db_name")
	if key == "" {
		key = "grom_cxt"
	}
	return DB.Set(key, c).WithContext(c)
}

// DbSurvive mysql survive
func DbSurvive(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		return err
	}
	return nil
}

const (
	callBackBeforeName = "opentracing:before"
	callBackAfterName  = "opentracing:after"
)

func Initialize(db *gorm.DB) (err error) {
	// 开始前 - 并不是都用相同的方法，可以自己自定义
	db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, func(db *gorm.DB) {})
	db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, func(db *gorm.DB) {})
	db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, func(db *gorm.DB) {})
	db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, func(db *gorm.DB) {})
	db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, func(db *gorm.DB) {})
	db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, func(db *gorm.DB) {})

	// 结束后 - 并不是都用相同的方法，可以自己自定义
	db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, func(db *gorm.DB) {})
	db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, func(db *gorm.DB) {})
	db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, func(db *gorm.DB) {})
	db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, func(db *gorm.DB) {})
	db.Callback().Row().After("gorm:row").Register(callBackAfterName, func(db *gorm.DB) {})
	db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, func(db *gorm.DB) {})
	return
}

// 记录主库写入时间，在查询的时候动态选择主库或从库
func registerReplicaCallbacks(db *gorm.DB) {
	db.Callback().Create().After("gorm:create").Register("record_write_time", recordWriteTime)
	db.Callback().Update().After("gorm:update").Register("record_write_time", recordWriteTime)
	db.Callback().Delete().After("gorm:delete").Register("record_write_time", recordWriteTime)
	db.Callback().Raw().After("gorm:raw").Register("record_write_time", recordRawWriteTime)
	db.Callback().Query().Before("gorm:query").Register("dynamic_read_write_clauses", dynamicReadWriteClauses)
	db.Callback().Row().Before("gorm:row").Register("dynamic_read_write_clauses", dynamicReadWriteClauses)
}

// 记录执行写入的时间
func recordWriteTime(db *gorm.DB) {
	ctx, ok := db.Statement.Context.(*gin.Context)
	if ok {
		ctx.Set("DBWriteTime", time.Now())
	}
}

// 记录执行写入的时间
func recordRawWriteTime(db *gorm.DB) {
	ctx, ok := db.Statement.Context.(*gin.Context)
	if ok {
		sql := db.Statement.SQL.String()
		if len(sql) >= 6 {
			prefix := strings.ToLower(sql[0:6])
			if prefix == "insert" || prefix == "update" || prefix == "delete" {
				ctx.Set("DBWriteTime", time.Now())
			}
		}
	}
}

// 动态选择使用主库还是从库
// 如果当前请求在 x 秒内刚执行过写入操作，则强制使用主库进行查询操作
func dynamicReadWriteClauses(db *gorm.DB) {
	ctx, ok := db.Statement.Context.(*gin.Context)
	if !ok {
		return
	}
	value, ok2 := ctx.Get("DBWriteTime")
	if !ok2 {
		return
	}
	writeTime, ok3 := value.(time.Time)
	if ok3 && time.Now().Sub(writeTime) <= dbWriteTime {
		db = db.Clauses(dbresolver.Write)
	}
}

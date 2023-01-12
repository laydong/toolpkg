package toolpkg

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/laydong/toolpkg/logx"
	"net/http"
	"time"
)

type AppConf struct {
	AppName       string        `json:"app_name"`       // 默认应用名称
	AppMode       string        `json:"app_mode"`       // 默认应用环境
	LogType       string        `json:"log_type"`       // 默认日志类型
	LogPath       string        `json:"log_path"`       // 默认文件目录
	ChildPath     string        `json:"child_path"`     // 默认子目录
	RotationSize  int           `json:"rotation_size"`  // 默认大小为32M
	RotationCount int           `json:"rotation_count"` // 默认不限制
	RotationTime  time.Duration `json:"rotation_time"`  // 默认每天轮转一次
	NoBuffWrite   bool          `json:"no_buff_write"`  // 不不开启无缓冲写入
	MaxAge        time.Duration `json:"max_age"`        // 默认保留90天
}

// InitLog 初始化日志服务
func InitLog(conf AppConf) {

	defaultConfig := logx.Config{
		AppName:       logx.DefaultAppName,
		AppMode:       logx.DefaultAppMode,
		LogType:       logx.DefaultLogType,
		LogPath:       logx.DefaultLogPath,
		ChildPath:     logx.DefaultChildPath,
		RotationSize:  int64(logx.DefaultRotationSize),
		RotationCount: uint(logx.DefaultRotationCount),
		NoBuffWrite:   logx.DefaultNoBuffWrite,
		RotationTime:  logx.DefaultRotationTime,
		MaxAge:        logx.DefaultMaxAge,
	}
	if conf.AppName != "" {
		defaultConfig.AppName = conf.AppName
	}
	if conf.AppMode != "" {
		defaultConfig.AppMode = conf.AppMode
	}
	if conf.LogType != "" {
		defaultConfig.LogType = conf.LogType
	}
	if conf.LogPath != "" {
		defaultConfig.LogPath = conf.LogPath
	}
	if conf.ChildPath != "" {
		defaultConfig.ChildPath = conf.ChildPath
	}
	if conf.RotationSize > 0 {
		defaultConfig.RotationSize = int64(conf.RotationSize)
	}

	if conf.RotationCount > 0 {
		defaultConfig.RotationCount = uint(conf.RotationCount)
	}
	if conf.RotationTime > 0 {
		defaultConfig.RotationTime = conf.RotationTime
	}
	if conf.NoBuffWrite != false {
		defaultConfig.NoBuffWrite = true
	}
	if conf.MaxAge > 0 {
		defaultConfig.MaxAge = conf.MaxAge
	}
	logx.InitLog(&defaultConfig)
}

// GetNewGinContext 获取新的上下文
func GetNewGinContext() *gin.Context {
	ctx := new(gin.Context)
	uid := uuid.New().String()
	ctx.Request = &http.Request{
		Header: make(map[string][]string),
	}
	ctx.Request.Header.Set(XtraceKey, uid)
	ctx.Request.Header.Set(RequestIdKey, uid)
	ctx.Set(RequestIdKey, uid)
	ctx.Set(XtraceKey, uid)
	return ctx
}

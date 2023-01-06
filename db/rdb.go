package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/laydong/toolpkg"
	"github.com/laydong/toolpkg/logx"
)

const (
	defaultRedisPoolMinIdle = 2 // 连接池空闲连接数量
)

// InitRdb 初始化redis
func InitRdb(addr, password string, num int) (db *redis.Client, err error) {
	options := redis.Options{
		Addr:     addr,
		Password: password,
		DB:       num,
	}
	if options.MinIdleConns == 0 {
		options.MinIdleConns = defaultRedisPoolMinIdle
	}
	db = redis.NewClient(&options)
	_, err = db.Ping(context.Background()).Result()
	if err == redis.Nil {
		logx.ErrorF(toolpkg.GetNewGinContext(), "Nil reply returned by Rdb when key does not exist", err.Error())
		return
	} else if err != nil {
		logx.ErrorF(toolpkg.GetNewGinContext(), "redis数据库链接错误", err.Error())
		return
	}
	err = RdbSurvive(db)
	if err != nil {
		logx.ErrorF(toolpkg.GetNewGinContext(), "redis数据库存活检测失败", err.Error())
	}
	return
}

// RdbSurvive redis存活检测
func RdbSurvive(db *redis.Client) error {
	err := db.Ping(context.Background()).Err()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

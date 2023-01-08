package db

import (
	"context"
	"github.com/go-redis/redis/v8"
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
		return
	} else if err != nil {
		return
	}
	err = RdbSurvive(db)
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

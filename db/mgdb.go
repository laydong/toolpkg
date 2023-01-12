package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// InitMongoDb mongodb 初始化
// dns string 示例 "mongodb://root:123456@127.0.0.1:27627"
// maxsize 连接池空闲连接数量
// timeOut 空闲时间
func InitMongoDb(dns string, maxsize, timeOut int) (db *mongo.Client, err error) {
	MdbOptions := options.Client().
		ApplyURI(dns).
		SetMaxPoolSize(uint64(maxsize)).
		SetMinPoolSize(uint64(timeOut))
	db, err = mongo.NewClient(MdbOptions)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = db.Connect(ctx)
	return
}

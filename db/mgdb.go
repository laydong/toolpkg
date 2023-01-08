package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

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

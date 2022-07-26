package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var Mdb *mongo.Client

func InitMongoDb(dns string, maxsize, timeOut int) {
	MdbOptions := options.Client().
		ApplyURI(dns).
		SetMaxPoolSize(uint64(maxsize)).
		SetMinPoolSize(uint64(timeOut))
	db, err := mongo.NewClient(MdbOptions)
	if err != nil {
		log.Printf("[app.gstore] mgdb error: %v", err.Error())
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = db.Connect(ctx)
	if err != nil {
		log.Printf("[app.gstore] mgdb error: %v", err.Error())
		panic(err)
	}
	log.Printf("[app.gstore] mongo success")
	Mdb = db
}

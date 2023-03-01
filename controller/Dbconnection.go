package controller

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var User_Collection *mongo.Collection
var Rdb *redis.Client

func DB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err1 := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	Iserror(err1)
	db := client.Database("Usersdb")
	User_Collection = db.Collection("Users")
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

}

package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var ctx context
var UserModel *mongo.Collection

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
}

var DatabaseCtx context.Context

func OpenMongoClient() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:1234@localhost:27018/?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&ssl=false"))
	if err != nil {
		log.Fatal(err)
	}

	DatabaseCtx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(DatabaseCtx)
	if err != nil {
		log.Fatal(err)
	}
	quickstart := client.Database("quickstart")
	UserModel = quickstart.Collection("users")
	return client
	// defer client.Disconnect(DatabaseCtx)
}

package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var ctx context
var UserModel *mongo.Collection
var SiteModel *mongo.Collection

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
}

type Site struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     string             `bson:"userId"`
	Name       string             `json:"name"`
	WebsiteUrl string             `json: websiteUrl`
	Deleted    bool               `json: Deleted`
}

var DatabaseCtx context.Context

func OpenMongoClient() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:1234@localhost:27018/?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&ssl=false"))
	if err != nil {
		log.Fatal(err)
	}

	// DatabaseCtx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	DatabaseCtx, _ = context.WithCancel(context.Background())
	err = client.Connect(DatabaseCtx)
	if err != nil {
		log.Fatal(err)
	}
	quickstart := client.Database("quickstart")
	UserModel = quickstart.Collection("users")
	SiteModel = quickstart.Collection("sites")
	return client
	// defer client.Disconnect(DatabaseCtx)
}

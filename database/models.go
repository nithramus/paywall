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
var OffreModel *mongo.Collection

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
}

type Site struct {
	WebOffreUrl string `json: webOffreUrl`
}

type Abonnement struct {
}

type Client struct {
}

type Offre struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID      string             `json:"userId" bson:"userId"`
	Name        string             `json:"name"`
	Deleted     bool               `json: deleted`
	Sites       []Site
	Abonnements []Abonnement
	Clients     []Client
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
	OffreModel = quickstart.Collection("offres")
	return client
	// defer client.Disconnect(DatabaseCtx)
}

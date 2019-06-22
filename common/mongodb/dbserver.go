package mongodb

import (
	"context"
	"time"

	_ "github.com/go-sql-driver/mysql" //Importing mysql connector for golang
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Instance : Singleton Instance
var Instance *mongo.Client
var Profile *mongo.Collection
var Travel *mongo.Collection
var Social *mongo.Collection

func openDB() {
	Instance, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = Instance.Connect(ctx)
	if err != nil {
		panic(err)
	}
	Profile = Instance.Database("glimpse").Collection("profile")
	Travel = Instance.Database("glimpse").Collection("travel")
	Social = Instance.Database("glimpse").Collection("social")
}

func init() {
	openDB()
}

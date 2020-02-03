package mongodb

import (
	"context"
	"fmt"
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

var Addr = "mongodb://mongo:27017"
var Username = "issacnitinmongod"
var Password = "iPhoneMyPh0ne!!"

func init() {
	openDB()
	go healthChecks()
}

func openDB() {
	fmt.Printf("openDB called")
	Instance, err := mongo.NewClient(options.Client().ApplyURI(Addr).SetAuth(options.Credential{
		AuthSource: "admin",
		Username:   Username,
		Password:   Password,
	}))
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = Instance.Connect(ctx)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	Profile = Instance.Database("glimpse").Collection("profile")
	Travel = Instance.Database("glimpse").Collection("travel")
	Social = Instance.Database("glimpse").Collection("social")
}

func healthChecks() {
	for true {
		if Instance == nil || Profile == nil || Travel == nil || Social == nil {
			openDB()
		}
		time.Sleep(10 * time.Second)
	}
}

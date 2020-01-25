package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client mongo.Client

func init() {
	lclient, err := mongo.NewClient(options.Client().ApplyURI(
		"mongodb+srv://dav-jordan:CSbtfu2021@mtg-go-tgpc5.mongodb.net/test?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = lclient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	databases, err := lclient.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(databases)
	client = *lclient
}

func InsertUser(email string, card string) {
	// insert test user into database
	collection := client.Database("mtg").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, bson.M{"email": email, "card": card})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.InsertedID)
}

func Close() {
	ctx, err := context.WithTimeout(context.Background(), 10*time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	client.Disconnect(ctx)
}

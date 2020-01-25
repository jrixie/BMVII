package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func InsertAlert(email string, card string, price float32, condition int,
	threshold float32) {
	// insert test user into database
	collection := client.Database("mtg").Collection("alerts")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, bson.M{"email": email, "cardID": card,
		"price": price, "condition": condition, "threshold": threshold})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.InsertedID)
}

func GetAlerts(email string) primitive.M {
	collection := client.Database("mtg").Collection("alerts")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, bson.M{"email": email})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	var result bson.M
	for cur.Next(ctx) {
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		// do something with result....
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return result
}

func Close() {
	ctx, err := context.WithTimeout(context.Background(), 10*time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	client.Disconnect(ctx)
}

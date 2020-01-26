package mongo

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client mongo.Client

// init initializes the package when loaded.
func init() {
	lClient, err := mongo.NewClient(options.Client().ApplyURI(
		"mongodb+srv://dav-jordan:CSbtfu2021@mtg-go-tgpc5.mongodb.net/test?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := lClient.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	Client = *lClient
}

// Close closes the connection to the DB.
func Close() {
	ctx, err := context.WithTimeout(context.Background(), 5*time.Second)
	if err != nil {
		log.Fatal("Could not create context")
	}

	if err := Client.Disconnect(ctx); err != nil {
		log.Fatal("Could not disconnect from DB")
	}
}

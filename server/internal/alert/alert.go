package alert

import (
	"boilermakevii/api/internal/mongo"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/MagicTheGathering/mtg-sdk-go"
	log "github.com/sirupsen/logrus"
)

// Trigger defines a data structure for user-defined
// price alerts, stored in the DB.
type Trigger struct {
	Email     string         `bson:"email"`
	CardID    mtg.CardId     `bson:"cardId"`
	Price     float32        `bson:"price"`
	Condition PriceCondition `bson:"condition"`
	Threshold float32        `bson:"threshold"`
}

// PriceCondition is a data type storing what condition to alert on.
type PriceCondition int8

const (
	GreaterThan PriceCondition = 0
	LessThan    PriceCondition = 1
)

// UpdateTriggers checks all entries from the DB, updates price information,
// and alerts users if their conditions have been met.
func UpdateTriggers() {
	var entries []Trigger

	entries, err := getAllEntries()
	if err != nil {
		log.Error(err)
	}

	for _, entry := range entries {
		// TODO: Update price

		if entry.hasMetCondition() {
			entry.alertUser()
		}
	}
}

// HasMetCondition returns true when a card price meets a
// user threshold, false otherwise.
func (t *Trigger) hasMetCondition() bool {
	switch t.Condition {
	case GreaterThan:
		return t.Price > t.Threshold
	case LessThan:
		return t.Price < t.Threshold
	default:
		return false
	}
}

// AlertUser alerts the user for a price alert.
func (t *Trigger) alertUser() {
	// TODO: Email user
}

// getAllEntries returns all Trigger entries from the DB.
func getAllEntries() ([]Trigger, error) {
	var entries []Trigger

	collection := mongo.Client.Database("mtg").Collection("alerts")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var entry Trigger
	for cur.Next(ctx) {
		if err := cur.Decode(&entry); err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	if err := cur.Close(ctx); err != nil {
		return nil, err
	}

	return entries, nil
}

func (t *Trigger) InsertTrigger() {
	collection := mongo.Client.Database("mtg").Collection("alerts")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	_, err := collection.InsertOne(ctx, t)
	if err != nil {
		log.Error(err)
	}
}

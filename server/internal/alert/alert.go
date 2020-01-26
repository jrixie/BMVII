package alert

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"boilermakevii/api/internal/mongo"
	"boilermakevii/api/internal/mtgjson"

	"github.com/MagicTheGathering/mtg-sdk-go"
	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

// Trigger defines a data structure for user-defined
// price alerts, stored in the DB.
type Trigger struct {
	Email     string         `bson:"email"`
	CardID    mtg.CardId     `bson:"cardId"`
	Price     float64        `bson:"price"`
	Condition PriceCondition `bson:"condition"`
	Threshold float64        `bson:"threshold"`
}

// PriceCondition is a data type storing what condition to alert on.
type PriceCondition int8

const (
	GreaterThan PriceCondition = 0
	LessThan    PriceCondition = 1
)

// init initializes the package when loaded.
func init() {
	gocron.Every(30).Seconds().Do(UpdateTriggers)
}

// UpdateTriggers checks all entries from the DB, updates price information,
// and alerts users if their conditions have been met.
func UpdateTriggers() {
	var entries []Trigger

	entries, err := getAllEntries()
	if err != nil {
		log.Error(err)
	}

	for _, entry := range entries {
		entry.Price = mtgjson.CardPrices[entry.CardID].Price

		if entry.hasMetCondition() {
			entry.alertUser()
		}
	}

	// TODO: Update in DB
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

// insertTrigger inserts a Trigger in the DB.
func (t *Trigger) insertTrigger() {
	collection := mongo.Client.Database("mtg").Collection("alerts")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	_, err := collection.InsertOne(ctx, t)
	if err != nil {
		log.Error(err)
	}
}

// Trigger received from frontend
type ClientTrigger struct {
	Email          string  `json:"email"`
	CardName       string  `json:"cardName"`
	PriceCondition int     `json:"priceCondition"`
	PriceThreshold float64 `json:"priceThreshold"`
}

// CreateTrigger handles a request to create a trigger in the DB.
func CreateTrigger(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)

	var trigger ClientTrigger
	log.Info(data)
	if err := json.Unmarshal(data, &trigger); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	arr, err := mtg.NewQuery().Where(mtg.CardName, trigger.CardName).All()
	if len(arr) == 0 && err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var card mtg.Card
	card = *arr[0]

	t := Trigger{
		Email:     trigger.Email,
		CardID:    card.Id,
		Price:     0,
		Condition: PriceCondition(trigger.PriceCondition),
		Threshold: trigger.PriceThreshold,
	}
	t.insertTrigger()
}

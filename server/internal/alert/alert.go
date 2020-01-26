package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/jordan-wright/email"

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
	Condition PriceCondition `bson:"condition"`
	Threshold float64        `bson:"threshold"`
}

// PriceCondition is a data type storing what condition to alert on.
type PriceCondition int8

const (
	GreaterThan PriceCondition = 0
	LessThan    PriceCondition = 1
)

var alertEmail string
var alertPass string

// init initializes the package when loaded.
func init() {
	alertEmail = os.Getenv("NOTIF_EMAIL")
	alertPass = os.Getenv("NOTIF_PASS")

	// TODO: Check if email was sent and don't send again
	gocron.Every(30).Seconds().Do(UpdateTriggers)
}

// UpdateTriggers checks all entries from the DB, updates price information,
// and alerts users if their conditions have been met.
func UpdateTriggers() {
	var entries []Trigger
	numEmails := 0

	entries, err := getAllEntries()
	if err != nil {
		log.Error(err)
	}

	for _, entry := range entries {
		price := mtgjson.CardPrices[entry.CardID].Price

		if entry.hasMetCondition(price) {
			entry.alertUser(price)
			entry.removeTrigger()
			numEmails++
		}
	}

	log.Info(fmt.Sprintf("UpdateTriggers: Sent %d emails.", numEmails))
}

// HasMetCondition returns true when a card price meets a
// user threshold, false otherwise.
func (t *Trigger) hasMetCondition(price float64) bool {
	switch t.Condition {
	case GreaterThan:
		return price > t.Threshold
	case LessThan:
		return price < t.Threshold
	default:
		return false
	}
}

// AlertUser alerts the user for a price alert.
func (t *Trigger) alertUser(currPrice float64) {
	cardData, err := t.CardID.Fetch()
	if err != nil {
		log.Error(err)
		return
	}

	var condition string
	switch t.Condition {
	case GreaterThan:
		condition = "rose"
	case LessThan:
		condition = "dropped"
	}

	emailBody := fmt.Sprintf("The price on %s has %s to your threshold of $%.2f.\nThe price is now $%.2f", cardData.Name, condition, t.Threshold, currPrice)

	e := email.NewEmail()
	e.From = fmt.Sprintf("MTGDrop <%s>", alertEmail)
	e.To = []string{t.Email}
	e.Subject = "MTGDrop: Price Alert"
	e.Text = []byte(emailBody)
	if err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", alertEmail, alertPass, "smtp.gmail.com")); err != nil {
		log.Error(err)
		log.Error("Failed to send email")
	}
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

// removeTrigger removes a Trigger from the DB.
func (t *Trigger) removeTrigger() {
	collection := mongo.Client.Database("mtg").Collection("alerts")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	_, err := collection.DeleteOne(ctx, t)
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
	if r.Method == "POST" {
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)

	var trigger ClientTrigger
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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))

	t := Trigger{
		Email:     trigger.Email,
		CardID:    card.Id,
		Condition: PriceCondition(trigger.PriceCondition),
		Threshold: trigger.PriceThreshold,
	}

	t.insertTrigger()
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

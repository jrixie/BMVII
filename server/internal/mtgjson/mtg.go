package mtgjson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/MagicTheGathering/mtg-sdk-go"
	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"
)

var client *http.Client

var CardPrices map[mtg.CardId]CardPrice

func init() {
	client = &http.Client{
		Timeout: 30 * time.Second,
	}

	gocron.Every(1).Day().At("05:00").Do(GetPrices)
}

type CardPrice struct {
	CardID mtg.CardId
	Price  float64
}

func GetPrices() {
	log.Info("GetPrices: Getting price data...")
	CardPrices = make(map[mtg.CardId]CardPrice, 0)

	res, err := client.Get("https://www.mtgjson.com/files/AllPrices.json")
	if err != nil {
		log.Error("GetPrices: Could not obtain pricing data")
		return
	}

	data, err := ioutil.ReadAll(res.Body)

	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		log.Fatal(err)
	}

	dateString := strings.Split(time.Now().String(), " ")[0]

	for id, obj := range rawData {
		cardType := obj.(map[string]interface{})["prices"].(map[string]interface{})["paper"]
		if cardType == nil {
			continue
		}

		price, _ := cardType.(map[string]interface{})[dateString].(float64)

		priceData := CardPrice{
			CardID: mtg.CardId(id),
			Price:  price,
		}

		CardPrices[priceData.CardID] = priceData
	}

	log.Info(fmt.Sprintf("GetPrices: Got pricing data for %d cards.", len(CardPrices)))
}

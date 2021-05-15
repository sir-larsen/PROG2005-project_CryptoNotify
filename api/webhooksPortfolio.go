package api

import (
	lib "CryptoNotify/coreLib"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
)

var portfolioWebhooks = make(map[string]lib.PortfolioWebhook)

//CheckPortfoliowebhooks - Function for checking the volume webhooks
func CheckPortfoliowebhooks() {
	iterat := Client.Collection(collectionPortfolio).Documents(Ctx)
	docSnaps, err := iterat.GetAll()
	if err != nil {
		fmt.Println(err)
		fmt.Errorf("SOMETHING WENT WRONG WITH Portfolio WEBHOOKS")
		return
	}

	for _, snap := range docSnaps {
		var webhook lib.PortfolioWebhook
		snap.DataTo(&webhook)
		webhook.WebhookID = snap.Ref.ID
		portfolioWebhooks[webhook.WebhookID] = webhook

		if webhook.GoRoutineExists == false {
			webhook.GoRoutineExists = true
			updatePortfolioWebhookVol(webhook)
			go updatePortfoliowebhook(webhook)
		}
	}
}

func updatePortfoliowebhook(webhook lib.PortfolioWebhook) {
	for {
		var Value float64
		for i, symbol := range webhook.Symbols {
			Value += webhook.Holdings[i] * lib.Cryptos[symbol].Price
		}
		webhook.CurrentValue = Value
		postPortfolioWebhook(webhook)
		time.Sleep(time.Duration(webhook.Timeout) * time.Second)
	}
}

//postVolumeWebhook - Function used for posting webhooks to the URL specified
func postPortfolioWebhook(webhook lib.PortfolioWebhook) {
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(webhook)
	http.Post(webhook.Url, "application/json", buffer)
	fmt.Println("Portfolio webhook with webhookID and symbol: ", webhook.WebhookID, ", ", webhook.Symbols, " has been sent")
	if err != nil {
		fmt.Println("ERROR IN POST OF VOLUME WEBHOOK", err)
	}
}

func updatePortfolioWebhookVol(webhook lib.PortfolioWebhook) error {
	_, err := Client.Collection(collectionPortfolio).Doc(webhook.WebhookID).Update(Ctx, []firestore.Update{
		{
			Path:  "GoRoutineExists",
			Value: webhook.GoRoutineExists,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

//VolumeWebhookReg - Intermediate function for adding webhooks regarding volume changes on the server
func PortfolioWebhookReg(w http.ResponseWriter, r *http.Request) {
	webhook, err := readPortHook(w, r)
	if err != nil {
		http.Error(w, "Something went wrong when adding webhook: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(webhook)
	AddPortfolioWebhook(webhook, w, r)
}

/*
{
    "url": 	"https://webhook.site/2417cf18-676c-4722-a97b-96739ffcc303",
    "phone_number": "+4793044522",
    "symbols": ["ETH", "BTC", "XRP"],
    "holdings": [2, 1, 3]
	"timeout": 30                                      --TIMEOUT IS IN SECONDS
}
*/
func readPortHook(w http.ResponseWriter, r *http.Request) (lib.PortfolioWebhook, error) {
	webhook := lib.PortfolioWebhook{}
	err := json.NewDecoder(r.Body).Decode(&webhook)
	if err != nil {
		return lib.PortfolioWebhook{}, err
	}

	//fmt.Println(len(webhook.Symbols))
	fmt.Println(len(webhook.Holdings))

	if len(webhook.Symbols) != len(webhook.Holdings) {
		return lib.PortfolioWebhook{}, errors.New("HOLDINGS AND SYMBOLS DOES NOT MATCH IN SIZE")
	}

	//Checking that the symbols exists in the structure
	for _, symbol := range webhook.Symbols {
		_, found := lib.Cryptos[symbol]
		//fmt.Println(found)
		if found == false {
			return lib.PortfolioWebhook{}, errors.New("One of the symbols provided did not exist")
		}
	}

	if webhook.Url == "" && webhook.Number == "" {
		return lib.PortfolioWebhook{}, errors.New("Neither url or number have been entered, provide at least one")
	}

	var startValue float64

	for i, symbol := range webhook.Symbols {
		startValue += webhook.Holdings[i] * lib.Cryptos[symbol].Price
		//fmt.Println(startValue)
	}
	webhook.StartValue = startValue
	webhook.CurrentValue = startValue
	webhook.GoRoutineExists = false
	fmt.Println(webhook)

	return webhook, nil
}

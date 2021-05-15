package api

import (
	lib "CryptoNotify/coreLib"
	"bytes"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"net/http"
	"errors"
)

var priceWebhooks = make(map[string]lib.PriceWebhook)

func CheckPriceWebhooks() {
	iterat := Client.Collection(collectionPrice).Documents(Ctx)
	docSnaps, err := iterat.GetAll()
	if err != nil {
		fmt.Println(err)
		fmt.Errorf("SOMETHING WENT WRONG WITH PRICE WEBHOOKS")
		return
	}

	for _, snap := range docSnaps {
		var webhook lib.PriceWebhook
		snap.DataTo(&webhook)
		webhook.WebhookID = snap.Ref.ID
		priceWebhooks[webhook.WebhookID] = webhook

		//updatePriceWebhook(webhook) //for later
	}
}

func updatePriceWebhook(webhook lib.PriceWebhook) {

	webhook.CurrentPrice = lib.Cryptos[webhook.Symbol].Price
	Triggered := false

	//Logic part of whether a price target has been hit
	if webhook.IsPriceIncrease == true {
		if webhook.CurrentPrice >= webhook.TargetPrice {
			Triggered = true
		}
	} else {
		if webhook.CurrentPrice <= webhook.TargetPrice {
			Triggered = true
		}
	}


	if Triggered == true{
		webhook.HasTriggered = true
		postPriceWebhook(webhook)
		// Delete webhook
	}else {

		err := updatePriceWebhookCurrent(webhook)
		if err != nil {
			fmt.Println(err)
			fmt.Println("WEBHOOK_VOLUME WITH FIREBASE_ID: ", webhook.WebhookID, " HAS GONE WRONG IN FIREBASE UPDATE OF CURRENT PRICE")
		}

	}
}



func postPriceWebhook(webhook lib.PriceWebhook) {
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(webhook)
	http.Post(webhook.Url, "application/json", buffer)
	fmt.Println("Price target webhook with webhookID and symbol: ", webhook.WebhookID, ", ", webhook.Symbol, " has been sent")
	if err != nil {
		fmt.Println("ERROR IN POST OF PRICE TARGET WEBHOOK", err)
	}
}




func updatePriceWebhookCurrent(webhook lib.PriceWebhook) error {
	_, err := Client.Collection(collectionPrice).Doc(webhook.WebhookID).Update(Ctx, []firestore.Update{
		{
			Path:  "CurrentPrice",
			Value: webhook.CurrentPrice,
		},
	})
	if err != nil {
		return err
	}
	return nil
}


func PriceWebhookReg(w http.ResponseWriter, r *http.Request) {
	webhook, err := readPriceHook(w, r)
	if err != nil {
		http.Error(w, "Something went wrong when adding webhook: "+err.Error(), http.StatusBadRequest)
		return
	}
	AddPriceWebhook(webhook, w, r)
}

//readPriceHook - Function for reading in the post request from webhook
/*Expected format for volume webhook body (example):
{
	"url": "webhook.site/something/something",	    //The URL you want the webhook to be posted to
	"phone_number": "+4797885707",					//Phone number you want to recieve messages to
	"symbol": "XRP",
	"target_price": 2                               //$2.0 is the price target

}
*/



func readPriceHook(w http.ResponseWriter, r *http.Request) (lib.PriceWebhook, error) {
	webhook := lib.PriceWebhook{}
	err := json.NewDecoder(r.Body).Decode(&webhook)
	if err != nil {
		return lib.PriceWebhook{}, err
	}

	//Checking webhook data for trash here. That the currency exists in the structure
	exist := false
	if _, ok := lib.Cryptos[webhook.Symbol]; ok {
		exist = true
	}
	if !exist { //If symbol doesn't exist
		return lib.PriceWebhook{}, errors.New("Currency is not tracked or doesn't exist")
	}

	//Expecting that the user is competent enough to enter a correct url or number
	if webhook.Url == "" && webhook.Number == "" {
		return lib.PriceWebhook{}, errors.New("Neither url or number have been entered, provide at least one")
	}

	//If come to this point, standard values will be inserted
	webhook.CurrentPrice = lib.Cryptos[webhook.Symbol].Price
	webhook.Name = lib.Cryptos[webhook.Symbol].Name
	webhook.HasTriggered = false
	webhook.StartPrice = lib.Cryptos[webhook.Symbol].Price
	if webhook.TargetPrice > webhook.StartPrice{
		webhook.IsPriceIncrease = true
	}else{
		webhook.IsPriceIncrease = false
	}
	

	return webhook, nil
}
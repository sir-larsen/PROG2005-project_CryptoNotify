package api

import (
	lib "CryptoNotify/coreLib"
	"bytes"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"net/http"
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


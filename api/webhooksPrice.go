package api

import (
	lib "CryptoNotify/coreLib"
	"fmt"
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

}

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

		//updatePriceWebhook(webhook)
	}
}

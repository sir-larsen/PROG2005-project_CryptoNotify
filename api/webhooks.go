package api

import (
	lib "CryptoNotify/coreLib"
	"encoding/json"
	"errors"
	"net/http"
)

//VolumeWebhookReg - Intermediate function for adding webhooks regarding volume changes on the server
func VolumeWebhookReg(w http.ResponseWriter, r *http.Request) {
	webhook, err := readVolHook(w, r)
	if err != nil {
		http.Error(w, "Something went wrong when adding webhook: "+err.Error(), http.StatusBadRequest)
		return
	}
	//ADD WEBHOOK TO FIREBASE HERE
	AddVolumeWebhook(webhook, w, r)
}

//readVolHook - Function for reading in the post request from webhook
/*Expected format for volume webhook body (example):
{
	"url": "webhook.site/something/something",	    //The URL you want the webhook to be posted to
	"phone_number": "+4795833037",					//Phone number you want to recieve messages to
	"symbol": "ETH",
	"percentage_threshold": 3                                 //3% increase in total currency volume
}
*/
func readVolHook(w http.ResponseWriter, r *http.Request) (lib.VolumeWebhook, error) {
	webhook := lib.VolumeWebhook{}
	err := json.NewDecoder(r.Body).Decode(&webhook)
	if err != nil {
		return lib.VolumeWebhook{}, err
	}

	//Checking webhook data for trash here. That the currency exists in the structure
	exist := false
	if _, ok := lib.Cryptos[webhook.Symbol]; ok {
		exist = true
	}
	if !exist { //If symbol doesn't exist
		return lib.VolumeWebhook{}, errors.New("Currency is not tracked or doesn't exist")
	}

	//Expecting that the user is competent enough to enter a correct url or number
	if webhook.Url == "" && webhook.Number == "" {
		return lib.VolumeWebhook{}, errors.New("Neither url or number have been entered, provide at least one")
	}

	//If come to this point, standard values will be inserted
	webhook.Id = lib.Cryptos[webhook.Symbol].Id
	webhook.Name = lib.Cryptos[webhook.Symbol].Name
	webhook.StartVol = lib.Cryptos[webhook.Symbol].Vol24
	webhook.CurrentVol = lib.Cryptos[webhook.Symbol].Vol24
	webhook.CurrentPercentage = 0
	webhook.HasTriggered = false

	return webhook, nil
}

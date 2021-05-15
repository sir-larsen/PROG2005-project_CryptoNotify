package api

import (
	lib "CryptoNotify/coreLib"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
)

var priceWebhooks = make(map[string]lib.PriceWebhook)

//CheckPriceWebhooks - Function for iterating through webhooks and checking for threshold reach
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

		updatePriceWebhook(webhook) //for later
	}
}

//updatePriceWebhook - Function for checking if price has been reached
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

	if Triggered == true {
		webhook.HasTriggered = true

		if webhook.Url != "" {
			postPriceWebhook(webhook)
		}

		if webhook.Number != "" {
			SendSmsFromPriceWebhook(webhook)
		}
		// Delete webhook
		DeletePriceWebhookInternal(webhook.WebhookID)
	} else {

		err := updatePriceWebhookCurrent(webhook)
		if err != nil {
			fmt.Println(err)
			fmt.Println("WEBHOOK_PRICE WITH FIREBASE_ID: ", webhook.WebhookID, " HAS GONE WRONG IN FIREBASE UPDATE OF CURRENT PRICE")
		}
	}
}

//postPriceWebhook - Function for POST of webhook to url
func postPriceWebhook(webhook lib.PriceWebhook) {
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(webhook)
	http.Post(webhook.Url, "application/json", buffer)
	fmt.Println("Price target webhook with webhookID and symbol: ", webhook.WebhookID, ", ", webhook.Symbol, " has been sent")
	if err != nil {
		fmt.Println("ERROR IN POST OF PRICE TARGET WEBHOOK", err)
	}
}

//updatePriceWebhook - Function for updating field in firebase
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

//WebhookPriceDel - Function dor user to delete a webhook
func WebhookPriceDel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") //Extracting the id
	if len(id) != 0 {
		DeletePriceWebhookFromAPI(w, r, id)
	} else {
		http.Error(w, "NO ID PROVIDED", http.StatusBadRequest)
	}
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
/*Expected format for target-price webhook body (example):
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
	if webhook.TargetPrice > webhook.StartPrice {
		webhook.IsPriceIncrease = true
	} else {
		webhook.IsPriceIncrease = false
	}

	return webhook, nil
}

//Function for rendering all the webhooks to the user
func AllPriceWebhooks(w http.ResponseWriter, r *http.Request) {
	var hooks []lib.PriceWebhook
	iter := Client.Collection(collectionPrice).Documents(Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		var hook lib.PriceWebhook
		doc.DataTo(&hook)
		hook.WebhookID = doc.Ref.ID

		hooks = append(hooks, hook)
	}
	w.Header().Add("content-type", "application/json")
	err := json.NewEncoder(w).Encode(hooks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//Fuction for getting webhook out to the browser
func GetPriceWebhook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") //Extracting the id
	dsnap, err := Client.Collection(collectionPrice).Doc(id).Get(Ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	var m lib.PriceWebhook
	dsnap.DataTo(&m)

	ref := Client.Collection(collectionPrice).Doc(id)
	m.WebhookID = ref.ID

	w.Header().Add("content-type", "application/json")
	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

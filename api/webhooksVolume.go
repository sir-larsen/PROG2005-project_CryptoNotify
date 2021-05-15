package api

import (
	lib "CryptoNotify/coreLib"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"google.golang.org/api/iterator"
)

var volumeWebhooks = make(map[string]lib.VolumeWebhook)

func CheckVolumeWebhooks() {
	iterat := Client.Collection(collectionVolume).Documents(Ctx)
	docSnaps, err := iterat.GetAll()
	if err != nil {
		fmt.Println(err)
		fmt.Errorf("SOMETHING WENT WRONG WITH VOLUME WEBHOOKS")
		return
	}

	for _, snap := range docSnaps {
		var webhook lib.VolumeWebhook
		snap.DataTo(&webhook)
		webhook.WebhookID = snap.Ref.ID
		volumeWebhooks[webhook.WebhookID] = webhook

		updateVolumeWebhook(webhook)
	}
}

func updateVolumeWebhook(webhook lib.VolumeWebhook) { //HUSK Å SKRIVE ENDRINGER TILBAKE TIL FIREBASE
	//DO ALL THE VOLUME STUFF CHECKS
	//IF TRIGGERED, SEND TO URL AND POSSIBLY PHONE NUMBA

	webhook.CurrentVol = lib.Cryptos[webhook.Symbol].Vol24 //Checking if the volume has reached the percentage threshold
	x := webhook.StartVol
	x /= 100
	x *= webhook.PercentThreshold

	//REGNE UT CURRENT PERCENTAGE OG LEGGE INN I WEBHOOK FØR SENDE TIL UPDATE
	x = webhook.StartVol //Figuring out the current percentage
	y := webhook.CurrentVol
	res := x / y
	res *= 100
	finalPercentage := 100 - res
	webhook.CurrentPercentage = finalPercentage //Updating the current percentage for neat tracking

	if webhook.CurrentVol >= webhook.StartVol+x { //Webhook has exceeded threshold, and is triggered
		webhook.HasTriggered = true

		//POST WEBHOOK
		postVolumeWebhook(webhook)

		//SMS NOTIFICATION   //TBA

		//DELETE WEBHOOK
		//DeleteVolumeWebhookInternal(webhook.WebhookID) HUSK Å KOMMENTER UT

	} else {
		//updateWebhookVolumeVol
		err := updateVolumeWebhookVol(webhook)
		if err != nil {
			fmt.Println(err)
			fmt.Println("WEBHOOK_VOLUME WITH FIREBASE_ID: ", webhook.WebhookID, " HAS GONE WRONG IN FIREBASE UPDATE OF CURRENTVOL")
		}

		//updateWebhookVolumePercentage
		err = updateVolumeWebhookPercentage(webhook)
		if err != nil {
			fmt.Println(err)
			fmt.Println("WEBHOOK_VOLUME WITH FIREBASE_ID: ", webhook.WebhookID, " HAS GONE WRONG IN FIREBASE UPDATE OF CURRENTPERCENTAGE")
		}
		//Send webhook to webhook site, not notification by sms, since then you can track changes in current percentage
		//postVolumeWebhook(webhook)
	}
}

//WebhookVolumeDel - Function dor user to delete a webhook
func WebhookVolumeDel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") //Extracting the id
	if len(id) != 0 {
		DeleteVolumeWebhookFromAPI(w, r, id)
	} else {
		http.Error(w, "NO ID PROVIDED", http.StatusBadRequest)
	}
}

//postVolumeWebhook - Function used for posting webhooks to the URL specified
func postVolumeWebhook(webhook lib.VolumeWebhook) {
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(webhook)
	http.Post(webhook.Url, "application/json", buffer)
	fmt.Println("Volume webhook with webhookID and symbol: ", webhook.WebhookID, ", ", webhook.Symbol, " has been sent")
	if err != nil {
		fmt.Println("ERROR IN POST OF VOLUME WEBHOOK", err)
	}
}

func updateVolumeWebhookVol(webhook lib.VolumeWebhook) error {
	_, err := Client.Collection(collectionVolume).Doc(webhook.WebhookID).Update(Ctx, []firestore.Update{
		{
			Path:  "CurrentVol",
			Value: webhook.CurrentVol,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func updateVolumeWebhookPercentage(webhook lib.VolumeWebhook) error {
	_, err := Client.Collection(collectionVolume).Doc(webhook.WebhookID).Update(Ctx, []firestore.Update{
		{
			Path:  "CurrentPercentage",
			Value: webhook.CurrentPercentage,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

//VolumeWebhookReg - Intermediate function for adding webhooks regarding volume changes on the server
func VolumeWebhookReg(w http.ResponseWriter, r *http.Request) {
	webhook, err := readVolHook(w, r)
	if err != nil {
		http.Error(w, "Something went wrong when adding webhook: "+err.Error(), http.StatusBadRequest)
		return
	}
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

//Function for rendering all the webhooks to the user
func AllVolumeWebhooks(w http.ResponseWriter, r *http.Request) {
	var hooks []lib.VolumeWebhook
	iter := Client.Collection(collectionVolume).Documents(Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		var hook lib.VolumeWebhook
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
func GetVolumeWebhook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") //Extracting the id
	dsnap, err := Client.Collection(collectionVolume).Doc(id).Get(Ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	var m lib.VolumeWebhook
	dsnap.DataTo(&m)

	ref := Client.Collection(collectionVolume).Doc(id)
	m.WebhookID = ref.ID

	w.Header().Add("content-type", "application/json")
	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

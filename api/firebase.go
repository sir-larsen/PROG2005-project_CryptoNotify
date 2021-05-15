package api

import (
	lib "CryptoNotify/coreLib"
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	//firestore "cloud.google.com/go/firestore/apiv1"
)

var collectionVolume = "webhooks_volume"
var collectionPrice = "webhooks_price"
var Ctx context.Context
var Client *firestore.Client
var projectID = "cloud-project-dd1b4"

func AddVolumeWebhook(webhook lib.VolumeWebhook, w http.ResponseWriter, r *http.Request) {

	ref, _, err := Client.Collection(collectionVolume).Add(Ctx, map[string]interface{}{
		"CurrentPercentage": webhook.CurrentPercentage,
		"CurrentVol":        webhook.CurrentVol,
		"HasTriggered":      webhook.HasTriggered,
		"Id":                webhook.Id,
		"Name":              webhook.Name,
		"Number":            webhook.Number,
		"PercentThreshold":  webhook.PercentThreshold,
		"StartVol":          webhook.StartVol,
		"Symbol":            webhook.Symbol,
		"Url":               webhook.Url,
	})
	if err != nil {
		http.Error(w, "Error when adding webhook "+webhook.Url, http.StatusBadRequest)
	} else {
		fmt.Println("Entry added to collection.")
		http.Error(w, ref.ID, http.StatusCreated) // Returns document ID
	}
}

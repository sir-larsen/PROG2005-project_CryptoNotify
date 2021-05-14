package api

import (
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
}

//readVolHook - Function for reading in the post request from webhook
/*
* Expected format for volume webhook body (example):
{
	"url": "webhook.site/something/something"	    //The URL you want the webhook to be posted to
	"phone_number": "+4795833037"					//Phone number you want to recieve messages to
	"symbol": "ETH"
	"percentage": 3                                 //3% increase in total currency volume
}
*
*/
func readVolHook(w http.ResponseWriter, r *http.Request) (VolumeWebhook, error) {
	webhook := VolumeWebhook{}

}

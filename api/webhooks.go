package api

import "net/http"

//VolumeWebhookReg - Intermediate function for adding webhooks regarding volume changes on the server
func VolumeWebhookReg(w http.ResponseWriter, r *http.Request) {
	webhook, err := readVolHook(w, r)
	if err != nil {
		http.Error(w, "Something went wrong when adding webhook: "+err.Error(), http.StatusBadRequest)
		return
	}
	//ADD WEBHOOK TO FIREBASE HERE
}

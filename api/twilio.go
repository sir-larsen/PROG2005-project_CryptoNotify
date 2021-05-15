package api

import (
	lib "CryptoNotify/coreLib"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var accountSid = "ACbe76999159d78a3662c88690bf3dbb8f" //Not optimal for security doing this
var authToken = "a331462d3a8a090916adc7d055ca5323"
var urlStr = "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
var twilioNum = "+13157534147"

func SendSmsFromVolumeWebhook(webhook lib.VolumeWebhook) {
	v := url.Values{}
	v.Set("To", webhook.Number)
	v.Set("From", twilioNum)
	threshold := fmt.Sprintf("%.2f", webhook.PercentThreshold)
	v.Set("Body", "Message from CryptoNotify! Your registered volume webhook hit its threshold of "+threshold+"%. Webhook will now be deleted")
	rb := *strings.NewReader(v.Encode())

	// Create client
	client := &http.Client{}

	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, _ := client.Do(req)
	fmt.Println(resp.Status)
}

func SendSmsFromPriceWebhook(webhook lib.PriceWebhook) {
	v := url.Values{}
	v.Set("To", webhook.Number)
	v.Set("From", twilioNum)
	pricePoint := fmt.Sprintf("%.2f", webhook.TargetPrice)
	v.Set("Body", "Message from CryptoNotify! Your registered price webhook hit its threshold of $"+pricePoint+" USD. Webhook will now be deleted")
	rb := *strings.NewReader(v.Encode())

	// Create client
	client := &http.Client{}

	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, _ := client.Do(req)
	fmt.Println(resp.Status)
}

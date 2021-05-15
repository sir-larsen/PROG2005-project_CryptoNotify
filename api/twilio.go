package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var accountSid = "ACbe76999159d78a3662c88690bf3dbb8f" //Not optimal for security doing this
var authToken = "a331462d3a8a090916adc7d055ca5323"
var urlStr = "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
var twilioNum = "+13157534147"

func SendMessage() {

	// Build out the data for our message
	v := url.Values{}
	v.Set("To", "+4793044522")
	v.Set("From", twilioNum)
	v.Set("Body", "Briefcase wanker")
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

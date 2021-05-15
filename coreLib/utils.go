package coreLib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var CryptoResp CryptoInfo                     //Global strut with all info regarding all the currencies
var Cryptos = make(map[string]CryptoInternal) //Map containing currencies for internal representation

//GetMock - Function for pulling the mocked info on currencies and putting them into in memory storage
func GetMock() {
	//"https://run.mocky.io/v3/ee4d32e9-1875-4f24-8e3e-1d9fb323bec0" old mock url
	resp, err := http.Get("https://9f878240-fe54-4229-ba8d-0ee03b66f0b9.mock.pstmn.io/cockandballs.com")
	if err != nil {
		fmt.Println(err, "SOMETHING WENT WRONG WHILE FETCHING DATA. RESTART SERVER!")
	}
	defer resp.Body.Close()
	//fmt.Println(resp.Body)
	err = json.NewDecoder(resp.Body).Decode(&CryptoResp)
	//fmt.Println(CryptoResp)
	fmt.Println("STATUS: MOCK PULLED")
	//fmt.Println(Cryptos)
}

func GetRealApi() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		log.Print(err)
	}

	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "100")
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "fa238227-46eb-4bc2-8e66-37c50f341fdb")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		//os.Exit(1)
	}
	//defer resp.Body.Close()
	fmt.Println(resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Error in reading response body from CMC")
	}
	json.Unmarshal(body, &CryptoResp)
	resp.Body.Close()
}

//UpdateInternalMap - Function for creating a map for internal representation of currencies
func UpdateInternalMap() {
	for _, cryptElemArr := range CryptoResp.Data {
		var currency CryptoInternal
		currency.Id = cryptElemArr.Id
		currency.Name = cryptElemArr.Name
		currency.Symbol = cryptElemArr.Symbol
		currency.MaxSupply = cryptElemArr.MaxSupply
		currency.CircSupply = cryptElemArr.CircSupply
		currency.TotSupply = cryptElemArr.TotSupply
		currency.Rank = cryptElemArr.Rank
		currency.Price = cryptElemArr.Quote.Usd.Price
		currency.Vol24 = cryptElemArr.Quote.Usd.Vol24
		currency.PercentChg24 = cryptElemArr.Quote.Usd.PercentChg24
		currency.PercentChg7d = cryptElemArr.Quote.Usd.PercentChg7d
		currency.MarketCap = cryptElemArr.Quote.Usd.MarketCap

		Cryptos[currency.Symbol] = currency
	}
	//fmt.Println(Cryptos)
	//fmt.Println(Cryptos["ADA"])
}

package coreLib

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var CryptoResp CryptoInfo                     //Global strut with all info regarding all the currencies
var Cryptos = make(map[string]CryptoInternal) //Map containing currencies for internal representation

//GetMock - Function for pulling the mocked info on currencies and putting them into in memory storage
func GetMock() {
	resp, err := http.Get("https://run.mocky.io/v3/ee4d32e9-1875-4f24-8e3e-1d9fb323bec0")
	if err != nil {
		fmt.Println(err, "SOMETHING WENT WRONG WHILE FETCHING DATA. RESTART SERVER!")
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&CryptoResp)
	fmt.Println("STATUS: MOCK PULLED")
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
}

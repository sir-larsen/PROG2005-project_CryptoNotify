package api

import (
	lib "CryptoNotify/coreLib"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

func CurrencyHandler(w http.ResponseWriter, r *http.Request) {
	currency := chi.URLParam(r, "currency_code")             //Currency symbol extracted from url
	outputStruct := lib.Cryptos[(strings.ToUpper(currency))] //Locating data from internal database
	//url := strings.Split(r.URL.String(), "/")
	fmt.Println(currency)

	if currency == "" {
		w.Header().Add("content-type", "application/json")
		err := json.NewEncoder(w).Encode( /*lib.CryptoResp*/ lib.Cryptos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("HHHHHHHH")
		}
	} else {

		if outputStruct.Name != "" { //If the currency is found

			w.Header().Add("content-type", "application/json")
			err := json.NewEncoder(w).Encode(outputStruct) //Display in json-format
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			// If symbol matching fails - Check for the full name

			currency = strings.Title(strings.ToLower(currency))
			var currencyName string
			for key, element := range lib.Cryptos {

				if element.Name == currency { //Iterate through map for matches on currency names.
					currencyName = key
				}
			}
			outputStruct2 := lib.Cryptos[(strings.ToUpper(currencyName))]

			if outputStruct2.Name != "" { //If there was a true match and no empty struct is created
				w.Header().Add("content-type", "application/json")
				err := json.NewEncoder(w).Encode(outputStruct2) //Display in json-format
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

			} else {
				//Currency not found after both symbol and name searching. Bad request
				http.Error(w, "Error: "+currency+" not found.", http.StatusBadRequest)
			}
		}
	}
}

func AllCurrencies(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("content-type", "application/json")
	err := json.NewEncoder(w).Encode( /*lib.CryptoResp*/ lib.Cryptos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
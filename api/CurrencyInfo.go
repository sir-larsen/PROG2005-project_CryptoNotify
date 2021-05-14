package api

import (
	lib "CryptoNotify/coreLib"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strings"
)

func CurrencyHandler(w http.ResponseWriter, r *http.Request) {
	currency := chi.URLParam(r, "currency_code") //Currency symbol extracted from url

	outputStruct := lib.Cryptos[(strings.ToUpper(currency))] //Locating data from internal database

	if outputStruct.Name != "" { //If the currency is found

		w.Header().Add("content-type", "application/json")
		err := json.NewEncoder(w).Encode(outputStruct) //Display in json-format
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	} else {
		var currencyName string
		for key, element := range lib.Cryptos {

			if element.Name == currency {

				currencyName = key
			}
		}
		outputStruct2 := lib.Cryptos[(strings.ToUpper(currencyName))]

		if outputStruct2.Name != "" {
			w.Header().Add("content-type", "application/json")
			err := json.NewEncoder(w).Encode(outputStruct2) //Display in json-format
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		}else {
			//Currency not found. Bad request
			http.Error(w, "Error: "+currency+" not found.", http.StatusBadRequest)
		}

	}
}
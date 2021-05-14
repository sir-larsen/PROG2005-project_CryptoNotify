package api

import (
	lib "CryptoNotify/coreLib"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

func CurrencyHandler(w http.ResponseWriter, r *http.Request) {
	currency := chi.URLParam(r, "currency_code")
	fmt.Print(currency)

	outputStruct := lib.Cryptos[currency]

	w.Header().Add("content-type", "application/json")
	err := json.NewEncoder(w).Encode(outputStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

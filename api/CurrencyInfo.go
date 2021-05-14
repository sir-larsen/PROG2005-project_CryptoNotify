package api

import (
	lib "CryptoNotify/coreLib"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strings"
)

func CurrencyHandler(w http.ResponseWriter, r *http.Request) {
	currency := chi.URLParam(r, "currency_code")

	outputStruct := lib.Cryptos[(strings.ToUpper(currency))]

	if outputStruct.Name != "" {

		w.Header().Add("content-type", "application/json")
		err := json.NewEncoder(w).Encode(outputStruct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	} else {
		http.Error(w, "Error: "+currency+" not found." ,http.StatusBadRequest)
	}
}

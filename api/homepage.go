package api

import (
	lib "CryptoNotify/coreLib"
	"encoding/json"
	"net/http"
)

// RootPage redirects to root
func RootPage(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Welcome to API")
	//fmt.Println(lib.Cock)

	w.Header().Add("content-type", "application/json")
	err := json.NewEncoder(w).Encode( /*lib.CryptoResp*/ lib.Cryptos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

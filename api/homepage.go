package api

import (
	"net/http"
)

// RootPage redirects to root
func RootPage(w http.ResponseWriter, r *http.Request) {

	htmlFilePath := "./api/index.html"
	http.ServeFile(w, r, htmlFilePath)

}

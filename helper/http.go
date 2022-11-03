package helper

import "net/http"

func E(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func CorsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Max-Age", "86400")
}

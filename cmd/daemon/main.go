package main

import (
	"bbrz/processor"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	reg := processor.NewRegistry()

	r.HandleFunc("/upload", reg.HandleProcessRequest).Methods(http.MethodPost)

	s := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}

	log.Fatalln(s.ListenAndServe())
}

package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alfreddobradi/go-bb-man/database"
	"github.com/alfreddobradi/go-bb-man/helper"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func ReplayListHandler(db database.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.GetReplayList()
		if err != nil {
			log.Printf("Failed to get replay list: %v", err)
			helper.E(w, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(rows); err != nil {
			log.Printf("Failed to encode response: %v", err)
			helper.E(w, http.StatusInternalServerError)
			return
		}
	}
}

func ReplayHandler(db database.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := uuid.Parse(vars["id"])
		if err != nil {
			log.Printf("Failed to parse ID %s: %v", vars["id"], err)
			helper.E(w, http.StatusInternalServerError)
			return
		}

		replay, err := db.GetReplay(id)
		if err != nil {
			log.Printf("Failed to get replay list: %v", err)
			helper.E(w, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(replay); err != nil {
			log.Printf("Failed to encode response: %v", err)
			helper.E(w, http.StatusInternalServerError)
			return
		}
	}
}

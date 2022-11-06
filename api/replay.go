package api

import (
	"encoding/json"
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
			logger.WithError(err).Error("Failed to get replay list")
			helper.E(w, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(rows); err != nil {
			logger.WithError(err).Error("Failed to encode response")
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
			logger.WithError(err).WithField("id", vars["id"]).Error("Failed to parse replay ID")
			helper.E(w, http.StatusInternalServerError)
			return
		}

		replay, err := db.GetReplay(id)
		if err != nil {
			logger.WithError(err).WithField("id", id).Error("Failed to get replay")
			helper.E(w, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(replay); err != nil {
			logger.WithError(err).Error("Failed to encode response")
			helper.E(w, http.StatusInternalServerError)
			return
		}
	}
}

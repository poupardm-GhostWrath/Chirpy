package main

import (
	"net/http"
	"encoding/json"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	const defaultEvent = "user.upgraded"

	type parameters struct {
		Event	string	`json:"event"`
		Data	struct {
			UserID	uuid.UUID	`json:"user_id"`
		} `json:"data"`
	}

	// Decode Request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Check Event
	if params.Event != defaultEvent {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Update User
	_, err = cfg.db.UpgradeUserByID(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
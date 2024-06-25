package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/lib/pq"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func isUniqueViolationError(err error) bool {
	var sqlErr *pq.Error
	if errors.As(err, &sqlErr) {
		// PostgreSQL unique constraint violation
		if sqlErr.Code == "23505" {
			return true
		}
	}
	return false
}

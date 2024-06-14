package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	type readyCheck struct {
		Status string `json:"status"`
	}
	respondWithJSON(w, http.StatusOK, readyCheck{
		Status: "ok",
	})
}

func handlerError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusOK, "Internal Server Error")
}

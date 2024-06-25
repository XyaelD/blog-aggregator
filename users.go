package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/XyaelD/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type createUserRequest struct {
		Name string
	}

	decoder := json.NewDecoder(r.Body)
	createRequest := createUserRequest{}
	err := decoder.Decode(&createRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt decode createUserRequest")
		return
	}

	if createRequest.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Missing name")
		return
	}

	ctx := r.Context()

	userUUID := uuid.New()
	currTime := time.Now().UTC()

	newUser := database.CreateUserParams{
		ID:        userUUID,
		CreatedAt: currTime,
		UpdatedAt: currTime,
		Name:      createRequest.Name,
	}

	createdUser, err := cfg.DB.CreateUser(ctx, newUser)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt create user")
		return
	}
	respondWithJSON(w, http.StatusCreated, dbUserToUser(createdUser))
}

func (cfg *apiConfig) handlerGetUserByApiKey(w http.ResponseWriter, r *http.Request, authedUser database.User) {
	respondWithJSON(w, http.StatusOK, dbUserToUser(authedUser))
}

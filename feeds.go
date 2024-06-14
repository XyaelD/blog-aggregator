package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/XyaelD/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, authedUser database.User) {
	type createFeedRequest struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	createRequest := createFeedRequest{}
	err := decoder.Decode(&createRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt decode createFeedRequest")
		return
	}

	if createRequest.Name == "" || createRequest.URL == "" {
		respondWithError(w, http.StatusBadRequest, "Request has empty values")
		return
	}

	ctx := r.Context()

	feedUUID := uuid.New()
	currTime := time.Now().UTC()

	newFeed := database.CreateFeedParams{
		ID:        feedUUID,
		CreatedAt: currTime,
		UpdatedAt: currTime,
		Name:      createRequest.Name,
		Url:       createRequest.URL,
		UserID:    authedUser.ID,
	}

	createdFeed, err := cfg.DB.CreateFeed(ctx, newFeed)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt create feed")
		return
	}

	feedFollowUUID := uuid.New()
	newFeedFollow := database.CreateFeedFollowParams{
		ID:        feedFollowUUID,
		CreatedAt: currTime,
		UpdatedAt: currTime,
		UserID:    createdFeed.UserID,
		FeedID:    createdFeed.ID,
	}
	createdFeedFollow, err := cfg.DB.CreateFeedFollow(ctx, newFeedFollow)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt create feed follow")
		return
	}

	type feedAndFollow struct {
		Feed       database.Feed       `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
	}

	results := feedAndFollow{
		Feed:       createdFeed,
		FeedFollow: createdFeedFollow,
	}
	respondWithJSON(w, http.StatusCreated, results)
}

func (cfg *apiConfig) handlerGetAllFeeds(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	allFeeds, err := cfg.DB.GetAllFeeds(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt retreive all feeds")
		return
	}
	respondWithJSON(w, http.StatusOK, allFeeds)
}

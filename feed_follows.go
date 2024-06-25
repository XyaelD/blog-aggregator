package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/XyaelD/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, authedUser database.User) {
	type createFeedFollowRequest struct {
		FeedId uuid.UUID //`json:"feed_id"`
	}

	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	createRequest := createFeedFollowRequest{}
	err := decoder.Decode(&createRequest)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt decode createFeedFollowRequest")
		return
	}

	// if createRequest.FeedId == "" {
	// 	respondWithError(w, http.StatusBadRequest, "request has empty values")
	// 	return
	// }

	// createRequestUUID, err := uuid.Parse(createRequest.FeedId)
	// if err != nil {
	// 	respondWithError(w, http.StatusBadRequest, "cannot parse request to UUID")
	// 	return
	// }

	// exists := false
	// allFeeds, err := cfg.DB.GetAllFeeds(ctx)
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "cannot check feeds in database")
	// 	return
	// }
	// for _, feed := range allFeeds {
	// 	if feed.ID == createRequestUUID {
	// 		exists = true
	// 		break
	// 	}
	// }
	// if !exists {
	// 	respondWithError(w, http.StatusBadRequest, "requested feed does not exist in the db")
	// 	return
	// }

	feedFollowUUID := uuid.New()
	currTime := time.Now().UTC()

	newFeedFollow := database.CreateFeedFollowParams{
		ID:        feedFollowUUID,
		CreatedAt: currTime,
		UpdatedAt: currTime,
		UserID:    authedUser.ID,
		FeedID:    createRequest.FeedId,
	}
	createdFeedFollow, err := cfg.DB.CreateFeedFollow(ctx, newFeedFollow)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt create feed follow")
		return
	}
	respondWithJSON(w, http.StatusOK, dbFeedFollowToFeedFollow(createdFeedFollow))
}

func (cfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, authedUser database.User) {
	feedFollowIDStr := r.PathValue("feedFollowID")
	deleteRequestUUID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not convert id")
		return
	}

	ctx := r.Context()

	deleteParams := database.DeleteFeedFollowParams{
		ID:     deleteRequestUUID,
		UserID: authedUser.ID,
	}

	err = cfg.DB.DeleteFeedFollow(ctx, deleteParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "feed follow does not exist")
		return
	}
	respondWithCode(w, http.StatusNoContent)
}

////////

func (cfg *apiConfig) handlerGetFeedFollowsForUser(w http.ResponseWriter, r *http.Request, authedUser database.User) {
	ctx := r.Context()
	feedFollows, err := cfg.DB.GetFeedFollowsForUser(ctx, authedUser.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot fetch feed follows from db")
		return
	}
	respondWithJSON(w, http.StatusFound, dbFeedFollowsToFeedFollows(feedFollows))
}

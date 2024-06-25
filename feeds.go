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
		Feed       Feed       `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
	}

	results := feedAndFollow{
		Feed:       dbFeedtoFeed(createdFeed),
		FeedFollow: dbFeedFollowToFeedFollow(createdFeedFollow),
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
	respondWithJSON(w, http.StatusOK, dbFeedsToFeeds(allFeeds))
}

//////////////

// func (cfg *apiConfig) fetchNextFeeds(limit int) ([]Feed, error) {
// 	ctx := context.Background()
// 	fetchedFeeds, err := cfg.DB.GetNextFeedsToFetch(ctx, int32(limit))
// 	if err != nil {
// 		return []Feed{}, err
// 	}
// 	modifiedFeeds := []Feed{}
// 	for _, feed := range fetchedFeeds {
// 		feedToAdd := dbFeedtoFeed(feed)
// 		modifiedFeeds = append(modifiedFeeds, feedToAdd)
// 	}
// 	return modifiedFeeds, nil
// }

// func (cfg *apiConfig) feedWorker(limit int) {
// 	ticker := time.NewTicker(60 * time.Second)
// 	defer ticker.Stop()

// 	for tickTime := range ticker.C {
// 		fmt.Println("Tick at", tickTime)
// 		fetchedFeeds, err := cfg.fetchNextFeeds(limit)
// 		if err != nil {
// 			fmt.Printf("Error fetching feeds: %v\n", err)
// 			return
// 		}

// 		var wg sync.WaitGroup
// 		for _, feed := range fetchedFeeds {
// 			wg.Add(1)
// 			go func(xmlURL string) {
// 				defer wg.Done()
// 				currRss, err := fetchDataFromFeed(xmlURL)
// 				if err != nil {
// 					log.Printf("Error fetching data from feed %s: %v\n", xmlURL, err)
// 					return
// 				}
// 				fmt.Printf("The current RSS title: %v\n", currRss.Channel.Title)
// 			}(feed.Url)
// 		}
// 		wg.Wait()
// 	}
// }

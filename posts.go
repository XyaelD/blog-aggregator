package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/XyaelD/blog-aggregator/internal/database"
)

var dateFormats = []string{
	time.RFC1123,          // "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC1123Z,         // "Mon, 02 Jan 2006 15:04:05 -0700"
	time.RFC3339,          // "2006-01-02T15:04:05Z07:00"
	"2006-01-02 15:04:05", // Custom format example
}

func parseDate(dateString string) time.Time {
	var parsedTime time.Time
	var err error

	for _, layout := range dateFormats {
		parsedTime, err = time.Parse(layout, dateString)
		if err == nil {
			return parsedTime
		}
	}
	return time.Time{}
}

func (cfg *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, authedUser database.User) {
	limit := r.URL.Query().Get("limit")
	ctx := r.Context()
	intLimit, err := strconv.ParseInt(limit, 10, 32)
	if err != nil {
		intLimit = 5
	}

	getPostParams := database.GetPostsByUserParams{
		UserID: authedUser.ID,
		Limit:  int32(intLimit),
	}

	userPosts, err := cfg.DB.GetPostsByUser(ctx, getPostParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot fetch posts")
		return
	}
	respondWithJSON(w, http.StatusFound, dbPostsToPosts(userPosts))
}

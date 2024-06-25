package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/XyaelD/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func scrapeWorker(cfg *apiConfig, db *database.Queries, concurrency int, timeBetweenRequests time.Duration) {
	log.Printf("Collecting feeds every %s on %v goroutines...", timeBetweenRequests, concurrency)
	ticker := time.NewTicker(timeBetweenRequests)

	ctx := context.Background()

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(ctx, int32(concurrency))
		if err != nil {
			log.Println("Couldn't get next feeds to fetch", err)
			continue
		}
		log.Printf("Found %v feeds to fetch!", len(feeds))

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(cfg, db, wg, feed)
		}
	}
}

func scrapeFeed(cfg *apiConfig, db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	ctx := context.Background()
	_, err := db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}
	feedData, err := fetchDataFromFeed(feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Items {
		// Create valid sql.NullString for description
		sqlString := sql.NullString{
			String: "",
			Valid:  false,
		}
		if item.Description != "" {
			sqlString = sql.NullString{
				String: item.Description,
				Valid:  true,
			}
		}
		// Create valid sql.NullTime for published_at
		sqlTime := sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
		emptyTime := time.Time{}
		parsedTime := parseDate(item.PubDate)
		if parsedTime != emptyTime {
			sqlTime = sql.NullTime{
				Time:  parsedTime,
				Valid: true,
			}
		}
		newPost := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sqlString,
			PublishedAt: sqlTime,
			FeedID:      feed.ID,
		}
		_, err := cfg.DB.CreatePost(ctx, newPost)
		if err != nil {
			if isUniqueViolationError(err) {
				log.Printf("Duplicate post for URL: %s, skipping...", newPost.Url)
				continue
			} else {
				log.Printf("Could not create post: %v, error: %v", newPost, err)
				continue
			}
		}
		log.Printf("Successfully created post with title: %s", newPost.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Items))
}

func fetchDataFromFeed(feedUrl string) (*RSS, error) {

	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := httpClient.Get(feedUrl)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	rssResults := RSS{}
	err = xml.Unmarshal(body, &rssResults)
	if err != nil {
		return &RSS{}, err
	}

	return &rssResults, nil
}

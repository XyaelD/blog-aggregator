package main

import (
	"database/sql"
	"time"

	"github.com/XyaelD/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func dbUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
	}
}

type Feed struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Url           string
	UserID        uuid.UUID
	LastFetchedAt *time.Time
}

func dbFeedtoFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:            dbFeed.ID,
		CreatedAt:     dbFeed.CreatedAt,
		UpdatedAt:     dbFeed.UpdatedAt,
		Name:          dbFeed.Name,
		Url:           dbFeed.Url,
		UserID:        dbFeed.UserID,
		LastFetchedAt: sqlNullTimeToTimePtr(dbFeed.LastFetchedAt),
	}
}

func sqlNullTimeToTimePtr(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func sqlNullStringToStringPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func dbFeedsToFeeds(dbFeeds []database.Feed) []Feed {
	result := make([]Feed, len(dbFeeds))
	for i, dbFeed := range dbFeeds {
		result[i] = dbFeedtoFeed(dbFeed)
	}
	return result
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func dbFeedFollowToFeedFollow(dbFeedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFollow.ID,
		CreatedAt: dbFeedFollow.CreatedAt,
		UpdatedAt: dbFeedFollow.UpdatedAt,
		UserID:    dbFeedFollow.UserID,
		FeedID:    dbFeedFollow.FeedID,
	}
}

func dbFeedFollowsToFeedFollows(dbFeedFollows []database.FeedFollow) []FeedFollow {
	result := make([]FeedFollow, len(dbFeedFollows))
	for i, feedFollow := range dbFeedFollows {
		result[i] = dbFeedFollowToFeedFollow(feedFollow)
	}
	return result
}

type Post struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `json:"title"`
	Url         string     `json:"url"`
	Description *string    `json:"description"`
	PublishedAt *time.Time `json:"published_at"`
	FeedId      uuid.UUID  `json:"feed_id"`
}

func dbPostToPost(dbPost database.Post) Post {
	return Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Url:         dbPost.Url,
		Description: sqlNullStringToStringPtr(dbPost.Description),
		PublishedAt: sqlNullTimeToTimePtr(dbPost.PublishedAt),
		FeedId:      dbPost.FeedID,
	}
}

func dbPostsToPosts(dbPosts []database.Post) []Post {
	result := make([]Post, len(dbPosts))
	for i, post := range dbPosts {
		result[i] = dbPostToPost(post)
	}
	return result
}

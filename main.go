package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/XyaelD/blog-aggregator/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}
	dbURL := os.Getenv("CONNECTION_STRING")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	apiConfig := &apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/healthz", handlerReadiness)
	mux.HandleFunc("GET /v1/err", handlerError)

	mux.HandleFunc("POST /v1/users", apiConfig.handlerCreateUser)
	mux.HandleFunc("GET /v1/users", apiConfig.middlewareAuth(apiConfig.handlerGetUserByApiKey))

	mux.HandleFunc("POST /v1/feeds", apiConfig.middlewareAuth(apiConfig.handlerCreateFeed))
	mux.HandleFunc("GET /v1/feeds", apiConfig.handlerGetAllFeeds)

	mux.HandleFunc("POST /v1/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerCreateFeedFollow))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiConfig.middlewareAuth(apiConfig.handlerDeleteFeedFollow))
	mux.HandleFunc("GET /v1/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerGetFeedFollowsForUser))

	mux.HandleFunc("GET /v1/posts", apiConfig.middlewareAuth(apiConfig.handlerGetPostsForUser))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	const concurrency = 1
	const requestInterval = 60 * time.Second

	go scrapeWorker(apiConfig, dbQueries, concurrency, requestInterval)

	log.Printf("Server started on port: %v\n", port)
	log.Fatal(server.ListenAndServe())
}

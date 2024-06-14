package main

import (
	"net/http"
	"strings"

	"github.com/XyaelD/blog-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiStr := r.Header.Get("Authorization")
		if apiStr == "" {
			respondWithError(w, http.StatusUnauthorized, "api key not present for this request")
			return
		}
		if !strings.HasPrefix(apiStr, "ApiKey ") {
			respondWithError(w, http.StatusUnauthorized, "invalid api key format")
			return
		}
		cleanApi := strings.TrimPrefix(apiStr, "ApiKey ")
		ctx := r.Context()

		selectedUser, err := cfg.DB.GetUserByApiKey(ctx, cleanApi)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "invalid api key")
			return
		}
		handler(w, r, selectedUser)
	}
}

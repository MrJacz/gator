package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mrjacz/gator/internal/database"
)

type FeedResponse struct {
	ID            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	URL           string     `json:"url"`
	UserID        uuid.UUID  `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_fetched_at,omitempty"`
}

type CreateFeedRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type FollowFeedRequest struct {
	FeedURL string `json:"feed_url"`
}

type FeedFollowResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func (s *Server) HandleCreateFeed(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CreateFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" || req.URL == "" {
		respondWithError(w, http.StatusBadRequest, "Name and URL are required")
		return
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      req.Name,
		Url:       req.URL,
		UserID:    userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create feed")
		return
	}

	// Automatically follow the feed
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
		FeedID:    feed.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to follow feed")
		return
	}

	respondWithJSON(w, http.StatusCreated, FeedResponse{
		ID:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		URL:           feed.Url,
		UserID:        feed.UserID,
		LastFetchedAt: nil,
	})
}

func (s *Server) HandleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch feeds")
		return
	}

	feedResponses := make([]FeedResponse, len(feeds))
	for i, feed := range feeds {
		var lastFetched *time.Time
		if feed.LastFetchedAt.Valid {
			lastFetched = &feed.LastFetchedAt.Time
		}

		feedResponses[i] = FeedResponse{
			ID:            feed.ID,
			CreatedAt:     feed.CreatedAt,
			UpdatedAt:     feed.UpdatedAt,
			Name:          feed.Name,
			URL:           feed.Url,
			UserID:        feed.UserID,
			LastFetchedAt: lastFetched,
		}
	}

	respondWithJSON(w, http.StatusOK, feedResponses)
}

func (s *Server) HandleFollowFeed(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req FollowFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.FeedURL == "" {
		respondWithError(w, http.StatusBadRequest, "Feed URL is required")
		return
	}

	feed, err := s.db.GetFeedByURL(context.Background(), req.FeedURL)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Feed not found")
		return
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
		FeedID:    feed.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to follow feed")
		return
	}

	respondWithJSON(w, http.StatusCreated, FeedFollowResponse{
		ID:        feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		UserID:    feedFollow.UserID,
		FeedID:    feedFollow.FeedID,
	})
}

func (s *Server) HandleUnfollowFeed(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	feedURL := vars["url"]

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Feed not found")
		return
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: userID,
		FeedID: feed.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to unfollow feed")
		return
	}

	respondWithJSON(w, http.StatusOK, SuccessResponse{Message: "Successfully unfollowed feed"})
}

func (s *Server) HandleGetFeedFollows(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch feed follows")
		return
	}

	responses := make([]FeedFollowResponse, len(feedFollows))
	for i, ff := range feedFollows {
		responses[i] = FeedFollowResponse{
			ID:        ff.ID,
			CreatedAt: ff.CreatedAt,
			UpdatedAt: ff.UpdatedAt,
			UserID:    ff.UserID,
			FeedID:    ff.FeedID,
		}
	}

	respondWithJSON(w, http.StatusOK, responses)
}

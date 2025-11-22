package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mrjacz/gator/internal/database"
)

type PostResponse struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	FeedID      uuid.UUID `json:"feed_id"`
}

type BookmarkResponse struct {
	ID        uuid.UUID    `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	UserID    uuid.UUID    `json:"user_id"`
	PostID    uuid.UUID    `json:"post_id"`
	Post      PostResponse `json:"post"`
}

type CreateBookmarkRequest struct {
	PostURL string `json:"post_url"`
}

func (s *Server) HandleGetPosts(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch posts")
		return
	}

	postResponses := make([]PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = PostResponse{
			ID:          post.ID,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			Title:       post.Title,
			URL:         post.Url,
			Description: post.Description,
			PublishedAt: post.PublishedAt,
			FeedID:      post.FeedID,
		}
	}

	respondWithJSON(w, http.StatusOK, postResponses)
}

func (s *Server) HandleSearchPosts(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		respondWithError(w, http.StatusBadRequest, "Search query is required")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	searchPattern := "%" + query + "%"

	posts, err := s.db.SearchPostsForUser(context.Background(), database.SearchPostsForUserParams{
		UserID: userID,
		Title:  searchPattern,
		Limit:  int32(limit),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to search posts")
		return
	}

	postResponses := make([]PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = PostResponse{
			ID:          post.ID,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			Title:       post.Title,
			URL:         post.Url,
			Description: post.Description,
			PublishedAt: post.PublishedAt,
			FeedID:      post.FeedID,
		}
	}

	respondWithJSON(w, http.StatusOK, postResponses)
}

func (s *Server) HandleCreateBookmark(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CreateBookmarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.PostURL == "" {
		respondWithError(w, http.StatusBadRequest, "Post URL is required")
		return
	}

	post, err := s.db.GetPostByURL(context.Background(), database.GetPostByURLParams{
		UserID: userID,
		Url:    req.PostURL,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	isBookmarked, err := s.db.IsPostBookmarked(context.Background(), database.IsPostBookmarkedParams{
		UserID: userID,
		PostID: post.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to check bookmark status")
		return
	}

	if isBookmarked {
		respondWithError(w, http.StatusConflict, "Post is already bookmarked")
		return
	}

	bookmark, err := s.db.CreateBookmark(context.Background(), database.CreateBookmarkParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
		PostID:    post.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create bookmark")
		return
	}

	respondWithJSON(w, http.StatusCreated, BookmarkResponse{
		ID:        bookmark.ID,
		CreatedAt: bookmark.CreatedAt,
		UpdatedAt: bookmark.UpdatedAt,
		UserID:    bookmark.UserID,
		PostID:    bookmark.PostID,
		Post: PostResponse{
			ID:          post.ID,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			Title:       post.Title,
			URL:         post.Url,
			Description: post.Description,
			PublishedAt: post.PublishedAt,
			FeedID:      post.FeedID,
		},
	})
}

func (s *Server) HandleDeleteBookmark(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	postURL := vars["url"]

	post, err := s.db.GetPostByURL(context.Background(), database.GetPostByURLParams{
		UserID: userID,
		Url:    postURL,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	err = s.db.DeleteBookmark(context.Background(), database.DeleteBookmarkParams{
		UserID: userID,
		PostID: post.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete bookmark")
		return
	}

	respondWithJSON(w, http.StatusOK, SuccessResponse{Message: "Bookmark deleted successfully"})
}

func (s *Server) HandleGetBookmarks(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	posts, err := s.db.GetBookmarksForUser(context.Background(), database.GetBookmarksForUserParams{
		UserID: userID,
		Limit:  int32(limit),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch bookmarks")
		return
	}

	postResponses := make([]PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = PostResponse{
			ID:          post.ID,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			Title:       post.Title,
			URL:         post.Url,
			Description: post.Description,
			PublishedAt: post.PublishedAt,
			FeedID:      post.FeedID,
		}
	}

	respondWithJSON(w, http.StatusOK, postResponses)
}

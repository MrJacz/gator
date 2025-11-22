package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}).Methods("GET")

	// Public routes
	r.HandleFunc("/api/users", s.HandleCreateUser).Methods("POST")
	r.HandleFunc("/api/users", s.HandleGetUsers).Methods("GET")
	r.HandleFunc("/api/users/{username}", s.HandleGetUser).Methods("GET")
	r.HandleFunc("/api/login", s.HandleLogin).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(AuthMiddleware)

	// Feed routes
	protected.HandleFunc("/feeds", s.HandleCreateFeed).Methods("POST")
	protected.HandleFunc("/feeds", s.HandleGetFeeds).Methods("GET")
	protected.HandleFunc("/feed_follows", s.HandleFollowFeed).Methods("POST")
	protected.HandleFunc("/feed_follows", s.HandleGetFeedFollows).Methods("GET")
	protected.HandleFunc("/feed_follows/{url}", s.HandleUnfollowFeed).Methods("DELETE")

	// Post routes
	protected.HandleFunc("/posts", s.HandleGetPosts).Methods("GET")
	protected.HandleFunc("/posts/search", s.HandleSearchPosts).Methods("GET")

	// Bookmark routes
	protected.HandleFunc("/bookmarks", s.HandleCreateBookmark).Methods("POST")
	protected.HandleFunc("/bookmarks", s.HandleGetBookmarks).Methods("GET")
	protected.HandleFunc("/bookmarks/{url}", s.HandleDeleteBookmark).Methods("DELETE")

	return r
}

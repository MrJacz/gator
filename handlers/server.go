package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mrjacz/gator/internal/api"
)

func Server(s *State, cmd Command) error {
	port := "8080"
	if len(cmd.Args) > 0 {
		port = cmd.Args[0]
	}

	// Check for JWT secret
	if os.Getenv("JWT_SECRET") == "" {
		log.Println("WARNING: JWT_SECRET not set, using default development secret")
		log.Println("Set JWT_SECRET environment variable for production use")
	}

	server := api.NewServer(s.DB)
	router := server.SetupRouter()

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting HTTP server on %s", addr)
	log.Printf("Health check: http://localhost%s/health", addr)
	log.Printf("API endpoints: http://localhost%s/api", addr)

	return http.ListenAndServe(addr, router)
}

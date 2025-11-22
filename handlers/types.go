package handlers

import (
	"github.com/mrjacz/gator/internal/config"
	"github.com/mrjacz/gator/internal/database"
)

// State holds the application State
type State struct {
	DB  *database.Queries
	Cfg *config.Config
}

// Command represents a CLI Command with its arguments
type Command struct {
	Name string
	Args []string
}

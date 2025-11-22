package main

import (
	"context"

	"github.com/mrjacz/gator/handlers"
	"github.com/mrjacz/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *handlers.State, cmd handlers.Command, user database.User) error) func(*handlers.State, handlers.Command) error {
	return func(s *handlers.State, cmd handlers.Command) error {
		user, err := s.DB.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}

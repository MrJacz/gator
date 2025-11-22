package main

import "github.com/mrjacz/gator/internal/database"


func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

}
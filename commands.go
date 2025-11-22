package main

import (
	"errors"

	"github.com/mrjacz/gator/handlers"
)

type commands struct {
	registeredCommands map[string]func(*handlers.State, handlers.Command) error
}

func (c *commands) register(name string, f func(*handlers.State, handlers.Command) error) {
	c.registeredCommands[name] = f
}

func (c *commands) run(s *handlers.State, cmd handlers.Command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}

package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/mrjacz/gator/handlers"
	"github.com/mrjacz/gator/internal/config"
	"github.com/mrjacz/gator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	programState := &handlers.State{
		DB:  dbQueries,
		Cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*handlers.State, handlers.Command) error),
	}
	cmds.register("register", handlers.Register)
	cmds.register("login", handlers.Login)
	cmds.register("reset", handlers.Reset)
	cmds.register("users", handlers.ListUsers)
	cmds.register("agg", handlers.Agg)
	cmds.register("server", handlers.Server)
	cmds.register("service", handlers.Service)
	cmds.register("addfeed", middlewareLoggedIn(handlers.AddFeed))
	cmds.register("feeds", handlers.ListFeeds)
	cmds.register("follow", middlewareLoggedIn(handlers.Follow))
	cmds.register("following", middlewareLoggedIn(handlers.ListFeedFollows))
	cmds.register("unfollow", middlewareLoggedIn(handlers.Unfollow))
	cmds.register("browse", middlewareLoggedIn(handlers.Browse))
	cmds.register("search", middlewareLoggedIn(handlers.Search))
	cmds.register("bookmark", middlewareLoggedIn(handlers.Bookmark))
	cmds.register("unbookmark", middlewareLoggedIn(handlers.Unbookmark))
	cmds.register("bookmarks", middlewareLoggedIn(handlers.ListBookmarks))
	cmds.register("tui", middlewareLoggedIn(handlers.TUI))

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, handlers.Command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mrjacz/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	sortBy := "date" // default sort by date
	var feedURL string

	// Parse arguments: browse [limit] [--sort=title|date] [--feed=url]
	for _, arg := range cmd.Args {
		if strings.HasPrefix(arg, "--sort=") {
			sortBy = strings.TrimPrefix(arg, "--sort=")
			if sortBy != "date" && sortBy != "title" {
				return fmt.Errorf("invalid sort option: %s (must be 'date' or 'title')", sortBy)
			}
		} else if strings.HasPrefix(arg, "--feed=") {
			feedURL = strings.TrimPrefix(arg, "--feed=")
		} else {
			// Assume it's the limit
			parsedLimit, err := strconv.Atoi(arg)
			if err != nil {
				return fmt.Errorf("invalid limit: %w", err)
			}
			limit = parsedLimit
		}
	}

	var posts []database.Post
	var err error

	// Fetch posts based on filters
	if feedURL != "" {
		posts, err = s.db.GetPostsForUserByFeed(context.Background(), database.GetPostsForUserByFeedParams{
			UserID: user.ID,
			Url:    feedURL,
			Limit:  int32(limit),
		})
	} else if sortBy == "title" {
		posts, err = s.db.GetPostsForUserSortedByTitle(context.Background(), database.GetPostsForUserSortedByTitleParams{
			UserID: user.ID,
			Limit:  int32(limit),
		})
	} else {
		posts, err = s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  int32(limit),
		})
	}

	if err != nil {
		return fmt.Errorf("couldn't get posts: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("No posts found.")
		return nil
	}

	fmt.Printf("Found %d posts for user %s", len(posts), user.Name)
	if feedURL != "" {
		fmt.Printf(" (filtered by feed: %s)", feedURL)
	}
	if sortBy == "title" {
		fmt.Printf(" (sorted by title)")
	} else {
		fmt.Printf(" (sorted by date)")
	}
	fmt.Println(":")

	for _, post := range posts {
		fmt.Printf("\n===================\n")
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("URL: %s\n", post.Url)
		fmt.Printf("Published: %s\n", post.PublishedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Description: %s\n", post.Description)
	}

	return nil
}

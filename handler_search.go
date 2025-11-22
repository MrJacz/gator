package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mrjacz/gator/internal/database"
)

func handlerSearch(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %s <search_term> [limit]", cmd.Name)
	}

	searchTerm := cmd.Args[0]
	limit := 10 // default limit for search results

	if len(cmd.Args) > 1 {
		// Try to parse limit from second argument
		var err error
		fmt.Sscanf(cmd.Args[1], "%d", &limit)
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	// Add wildcards for ILIKE pattern matching
	searchPattern := "%" + searchTerm + "%"

	posts, err := s.db.SearchPostsForUser(context.Background(), database.SearchPostsForUserParams{
		UserID: user.ID,
		Title:  searchPattern,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't search posts: %w", err)
	}

	if len(posts) == 0 {
		fmt.Printf("No posts found matching '%s'.\n", searchTerm)
		return nil
	}

	fmt.Printf("Found %d post(s) matching '%s':\n", len(posts), searchTerm)

	for _, post := range posts {
		fmt.Printf("\n===================\n")
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("URL: %s\n", post.Url)
		fmt.Printf("Published: %s\n", post.PublishedAt.Format("2006-01-02 15:04:05"))

		// Highlight search term in description (simple approach)
		description := post.Description
		if len(description) > 200 {
			// Find the search term and show context around it
			lowerDesc := strings.ToLower(description)
			lowerTerm := strings.ToLower(searchTerm)
			idx := strings.Index(lowerDesc, lowerTerm)

			if idx != -1 {
				// Show 100 chars before and after the match
				start := idx - 100
				if start < 0 {
					start = 0
				}
				end := idx + len(searchTerm) + 100
				if end > len(description) {
					end = len(description)
				}

				snippet := description[start:end]
				if start > 0 {
					snippet = "..." + snippet
				}
				if end < len(description) {
					snippet = snippet + "..."
				}
				description = snippet
			} else {
				// If not found in lowercase comparison, just truncate
				description = description[:200] + "..."
			}
		}

		fmt.Printf("Description: %s\n", description)
	}

	return nil
}

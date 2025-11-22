package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mrjacz/gator/internal/database"
)

func handlerBookmark(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <post_url>", cmd.Name)
	}

	postURL := cmd.Args[0]

	// Find the post by URL
	post, err := s.db.GetPostByURL(context.Background(), database.GetPostByURLParams{
		UserID: user.ID,
		Url:    postURL,
	})
	if err != nil {
		return fmt.Errorf("post not found with URL: %s", postURL)
	}

	// Check if already bookmarked
	isBookmarked, err := s.db.IsPostBookmarked(context.Background(), database.IsPostBookmarkedParams{
		UserID: user.ID,
		PostID: post.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't check bookmark status: %w", err)
	}

	if isBookmarked {
		return fmt.Errorf("post is already bookmarked")
	}

	// Create bookmark
	_, err = s.db.CreateBookmark(context.Background(), database.CreateBookmarkParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		PostID:    post.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create bookmark: %w", err)
	}

	fmt.Printf("Bookmarked: %s\n", post.Title)
	return nil
}

func handlerUnbookmark(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <post_url>", cmd.Name)
	}

	postURL := cmd.Args[0]

	// Find the post by URL
	post, err := s.db.GetPostByURL(context.Background(), database.GetPostByURLParams{
		UserID: user.ID,
		Url:    postURL,
	})
	if err != nil {
		return fmt.Errorf("post not found with URL: %s", postURL)
	}

	// Delete bookmark
	err = s.db.DeleteBookmark(context.Background(), database.DeleteBookmarkParams{
		UserID: user.ID,
		PostID: post.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't remove bookmark: %w", err)
	}

	fmt.Printf("Removed bookmark: %s\n", post.Title)
	return nil
}

func handlerListBookmarks(s *state, cmd command, user database.User) error {
	limit := 10 // default limit

	if len(cmd.Args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = parsedLimit
	}

	posts, err := s.db.GetBookmarksForUser(context.Background(), database.GetBookmarksForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get bookmarks: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("No bookmarks found.")
		return nil
	}

	fmt.Printf("Found %d bookmark(s):\n", len(posts))

	for _, post := range posts {
		fmt.Printf("\n===================\n")
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("URL: %s\n", post.Url)
		fmt.Printf("Published: %s\n", post.PublishedAt.Format("2006-01-02 15:04:05"))
		if len(post.Description) > 200 {
			fmt.Printf("Description: %s...\n", post.Description[:200])
		} else {
			fmt.Printf("Description: %s\n", post.Description)
		}
	}

	return nil
}

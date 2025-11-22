package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mrjacz/gator/internal/database"
	"github.com/mrjacz/gator/internal/rss"
)

func Agg(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %v <time_between_reqs> [--concurrency=N]", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	concurrency := 1 // default: fetch 1 feed at a time
	for _, arg := range cmd.Args[1:] {
		if strings.HasPrefix(arg, "--concurrency=") {
			concStr := strings.TrimPrefix(arg, "--concurrency=")
			parsedConc, err := strconv.Atoi(concStr)
			if err != nil {
				return fmt.Errorf("invalid concurrency value: %w", err)
			}
			if parsedConc < 1 {
				return fmt.Errorf("concurrency must be >= 1")
			}
			concurrency = parsedConc
		}
	}

	log.Printf("Collecting feeds every %s with concurrency %d...", timeBetweenRequests, concurrency)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s, concurrency)
	}
}

func scrapeFeeds(s *State, concurrency int) {
	feeds, err := s.DB.GetNextFeedsToFetch(context.Background(), int32(concurrency))
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}

	if len(feeds) == 0 {
		log.Println("No feeds to fetch")
		return
	}

	log.Printf("Found %d feed(s) to fetch!", len(feeds))

	var wg sync.WaitGroup
	for _, feed := range feeds {
		wg.Add(1)
		go func(f database.Feed) {
			defer wg.Done()
			scrapeFeed(s.DB, f)
		}(feed)
	}

	wg.Wait()
	log.Println("Finished fetching all feeds in this batch")
}

func parsePublishedAt(pubDate string) (time.Time, error) {
	// RSS feeds can use various date formats, try common ones
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		"2006-01-02T15:04:05Z07:00", // ISO 8601
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		t, err := time.Parse(format, pubDate)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", pubDate)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}

	for _, item := range feedData.Channel.Item {
		publishedAt, err := parsePublishedAt(item.PubDate)
		if err != nil {
			log.Printf("Couldn't parse published date for post '%s': %v", item.Title, err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			// Check if it's a duplicate URL error
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				// Ignore duplicate posts
				continue
			}
			// Log other errors
			log.Printf("Couldn't create post '%s': %v", item.Title, err)
			continue
		}

		log.Printf("Post created: %s", item.Title)
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}

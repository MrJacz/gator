<p align="center">
    <a href="https://www.boot.dev?bannerlord=itsmejacz">
        <img src="https://github.com/bootdotdev/bootdev/assets/4583705/7a1184f1-bb43-45fa-a363-f18f8309056f" />
    </a>
</p>

# Boot.dev Guided Project

This repository contains my solution to a guided project from [Boot.dev](https://www.boot.dev?bannerlord=jacz).

## Gator - RSS Feed Aggregator CLI

Gator is a command-line tool for aggregating and managing RSS feeds. It allows you to follow your favorite blogs and news sources, automatically fetch new posts, and browse them from your terminal.

## Prerequisites

Before running Gator, you need to have the following installed:

- **Go** (version 1.21 or higher) - [Download Go](https://go.dev/dl/)
- **PostgreSQL** - [Download PostgreSQL](https://www.postgresql.org/download/)

## Installation

### 1. Install the Gator CLI

You can install Gator using Go's built-in package manager:

```bash
go install github.com/mrjacz/gator@latest
```

This will compile the binary and place it in your `$GOPATH/bin` directory. Make sure this directory is in your system's `PATH`.

Alternatively, you can build it locally:

```bash
git clone https://github.com/mrjacz/gator.git
cd gator
go build -o gator
```

### 2. Set Up the Database

Create a PostgreSQL database for Gator:

```sql
CREATE DATABASE gator;
```

Run the database migrations using goose:

```bash
cd sql/schema
goose postgres "postgres://username:password@localhost:5432/gator?sslmode=disable" up
```

Replace `username` and `password` with your PostgreSQL credentials.

### 3. Configure Gator

Gator uses a configuration file located at `~/.gatorconfig.json`. Create this file with the following structure:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Replace `username` and `password` with your PostgreSQL credentials. The `current_user_name` will be set automatically when you register or login.

## Usage

### User Management

**Register a new user:**
```bash
gator register <username>
```

**Login as a user:**
```bash
gator login <username>
```

**List all users:**
```bash
gator users
```

**Reset the database (delete all users):**
```bash
gator reset
```

### Feed Management

**Add a new feed:**
```bash
gator addfeed <feed_name> <feed_url>
```

Example:
```bash
gator addfeed "Hacker News" "https://hnrss.org/newest"
gator addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"
```

**List all feeds:**
```bash
gator feeds
```

**Follow a feed:**
```bash
gator follow <feed_url>
```

**Unfollow a feed:**
```bash
gator unfollow <feed_url>
```

**List feeds you're following:**
```bash
gator following
```

### Aggregating Posts

**Start the aggregator (fetch posts from feeds):**
```bash
gator agg <time_between_requests> [--concurrency=N]
```

Examples:
```bash
gator agg 1m                    # Fetch 1 feed every 1 minute (sequential)
gator agg 30s                   # Fetch 1 feed every 30 seconds
gator agg 10s --concurrency=5   # Fetch 5 feeds concurrently every 10 seconds
gator agg 1m --concurrency=10   # Fetch 10 feeds concurrently every 1 minute
```

The aggregator will continuously fetch posts from feeds. Use `--concurrency` to fetch multiple feeds simultaneously for faster updates. Press `Ctrl+C` to stop it.

### Browse Posts

**View recent posts:**
```bash
gator browse [limit] [--sort=date|title] [--feed=feed_url] [--page=N]
```

Examples:
```bash
gator browse                                      # Shows 2 most recent posts (default)
gator browse 10                                   # Shows 10 most recent posts
gator browse 5 --sort=title                       # Shows 5 posts sorted alphabetically by title
gator browse 10 --feed="https://hnrss.org/newest" # Shows 10 posts from specific feed
gator browse 5 --page=2                           # Shows posts 6-10 (page 2 with default limit of 2 becomes 5)
gator browse 10 --page=3                          # Shows posts 21-30
gator browse 5 --sort=title --feed="https://blog.boot.dev/index.xml"  # Combine filters
```

Posts are displayed with their title, URL, publication date, and description.

### Search Posts

**Search for posts by keyword:**
```bash
gator search <search_term> [limit]
```

Examples:
```bash
gator search "golang"           # Search for posts containing "golang" (default 10 results)
gator search "Python" 20        # Search for posts containing "Python", show 20 results
gator search "API design" 5     # Search for posts about "API design", show 5 results
```

The search command performs case-insensitive searches across both post titles and descriptions, and displays matching posts with context snippets.

### Bookmark Posts

**Save posts for later reading:**
```bash
gator bookmark <post_url>         # Bookmark a post
gator unbookmark <post_url>       # Remove a bookmark
gator bookmarks [limit]           # List your bookmarks
```

Examples:
```bash
gator bookmark "https://blog.boot.dev/golang/benefits-of-go/"   # Bookmark a post
gator bookmarks                  # List all bookmarks (default 10)
gator bookmarks 20               # List 20 most recent bookmarks
gator unbookmark "https://blog.boot.dev/golang/benefits-of-go/" # Remove bookmark
```

Bookmarks are user-specific and persist across sessions.

## Example Workflow

```bash
# Register and login
gator register alice
gator login alice

# Add and follow some feeds
gator addfeed "Hacker News" "https://hnrss.org/newest"
gator addfeed "The Pragmatic Engineer" "https://blog.pragmaticengineer.com/rss/"

# Start aggregating (in a separate terminal or background)
gator agg 5m

# Browse the latest posts
gator browse 5
```

## Features

- ✅ Multi-user support with simple authentication
- ✅ Follow multiple RSS feeds
- ✅ Automatic feed aggregation with configurable intervals and concurrent fetching
- ✅ Duplicate post detection
- ✅ Robust date parsing for various RSS formats
- ✅ Browse posts with sorting (by date or title), filtering (by feed), and pagination
- ✅ Full-text search across post titles and descriptions
- ✅ Bookmark posts for later reading
- ✅ PostgreSQL database for persistent storage

## Technologies Used

- **Go** - Primary programming language
- **PostgreSQL** - Database
- **sqlc** - Type-safe SQL query generation
- **goose** - Database migrations
- **RSS/Atom** - Feed parsing

## License

This project is part of the Boot.dev curriculum and is for educational purposes.

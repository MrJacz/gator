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
gator agg <time_between_requests>
```

Example:
```bash
gator agg 1m  # Fetch feeds every 1 minute
gator agg 30s # Fetch feeds every 30 seconds
```

The aggregator will continuously fetch posts from all feeds. Press `Ctrl+C` to stop it.

### Browse Posts

**View recent posts:**
```bash
gator browse [limit] [--sort=date|title] [--feed=feed_url]
```

Examples:
```bash
gator browse                                      # Shows 2 most recent posts (default)
gator browse 10                                   # Shows 10 most recent posts
gator browse 5 --sort=title                       # Shows 5 posts sorted alphabetically by title
gator browse 10 --feed="https://hnrss.org/newest" # Shows 10 posts from specific feed
gator browse 5 --sort=title --feed="https://blog.boot.dev/index.xml"  # Combine filters
```

Posts are displayed with their title, URL, publication date, and description.

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
- ✅ Automatic feed aggregation with configurable intervals
- ✅ Duplicate post detection
- ✅ Robust date parsing for various RSS formats
- ✅ Browse posts with sorting (by date or title) and filtering (by feed)
- ✅ PostgreSQL database for persistent storage

## Technologies Used

- **Go** - Primary programming language
- **PostgreSQL** - Database
- **sqlc** - Type-safe SQL query generation
- **goose** - Database migrations
- **RSS/Atom** - Feed parsing

## License

This project is part of the Boot.dev curriculum and is for educational purposes.

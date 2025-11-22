-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.* FROM posts
JOIN feeds ON posts.feed_id = feeds.id
WHERE feeds.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;

-- name: GetPostsForUserByFeed :many
SELECT posts.* FROM posts
JOIN feeds ON posts.feed_id = feeds.id
WHERE feeds.user_id = $1 AND feeds.url = $2
ORDER BY posts.published_at DESC
LIMIT $3;

-- name: GetPostsForUserSortedByTitle :many
SELECT posts.* FROM posts
JOIN feeds ON posts.feed_id = feeds.id
WHERE feeds.user_id = $1
ORDER BY posts.title ASC
LIMIT $2;


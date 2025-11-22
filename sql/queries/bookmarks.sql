-- name: CreateBookmark :one
INSERT INTO bookmarks (id, created_at, updated_at, user_id, post_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteBookmark :exec
DELETE FROM bookmarks
WHERE user_id = $1 AND post_id = $2;

-- name: GetBookmarksForUser :many
SELECT posts.* FROM posts
JOIN bookmarks ON posts.id = bookmarks.post_id
WHERE bookmarks.user_id = $1
ORDER BY bookmarks.created_at DESC
LIMIT $2;

-- name: IsPostBookmarked :one
SELECT EXISTS(
    SELECT 1 FROM bookmarks
    WHERE user_id = $1 AND post_id = $2
);

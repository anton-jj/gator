-- name: CreatePost :one
INSERT INTO posts(id, created_at, updated_at, title, url, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
	$5,
	$6
)
RETURNING *;
-- name: GetPosts :many
	SELECT * FROM posts WHERE feed_id = $1
	ORDER BY created_at ASC LIMIT $2;

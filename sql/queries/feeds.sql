-- name: CreateFeeds :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
	$5,
	$6
)
RETURNING *;
-- name: GetFeed :one
	SELECT * FROM users WHERE id = $1;
-- name: GetFeeds :many
	SELECT * FROM feeds;
-- name: MarkFetchedFeed :exec
    INSERT INTO feeds (updated_at, last_fetched_at) SELECT $1, $2 WHERE EXISTS (SELECT 1 FROM feeds WHERE feeds.id = $3);
-- name: GetNextFeedToFetch :one
SELECT *
    FROM feeds 
    ORDER BY last_fetched_at ASC NULLS FIRST LIMIT 1 ;

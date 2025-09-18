-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
	INSERT INTO feed_follows (id, user_id, feed_id, updated_at, created_at)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at, user_id, feed_id
)
SELECT
	iif.*,
	f.name AS feed_name,
	u.name AS user_name
	FROM inserted_feed_follow AS iif
	JOIN users AS u ON u.id = iif.user_id
	JOIN feeds AS f ON f.id = iif.feed_id;

-- name: FindFeedByUrl :one
SELECT id, name, url FROM feeds WHERE url = $1;

-- name: GetFeedFollowsForUser :many
SELECT ff.id, ff.user_id, ff.feed_id, ff.created_at, ff.updated_at, f.name AS feed_name, u.name AS user_name FROM feed_follows AS ff
INNER JOIN feeds AS f on ff.feed_id = f.id
INNER JOIN users AS u on ff.user_id = u.id
WHERE u.id = $1;


-- name: RemoveFeedFollowsRecord :exec
DELETE FROM feed_follows 
WHERE user_id = $1 AND feed_id = $2;

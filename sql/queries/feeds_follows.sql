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



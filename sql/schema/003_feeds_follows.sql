-- +goose Up
CREATE TABLE feed_follows(
	id UUID NOT NULL PRIMARY KEY, 
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	user_id UUID NOT NULL, 
	feed_id UUID NOT NULL, 
	CONSTRAINT fk_feed_follows_user 
	FOREIGN KEY (user_id)
	REFERENCES users(id)
	ON DELETE CASCADE,

	CONSTRAINT fk_feed_follows_feed
	FOREIGN KEY (feed_id)
	REFERENCES feeds(id)
	ON DELETE CASCADE,
	CONSTRAINT uq_feed_follows_user_feed UNIQUE (user_id, feed_id)
);
-- +goose Down
DROP TABLE feed_follows;

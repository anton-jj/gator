-- +goose Up
CREATE TABLE posts(
	id UUID PRIMARY KEY, 
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	title TEXT NOT NULL,
	url TEXT UNIQUE NOT NULL,
	description TEXT,
	published_at DATE,
	feed_id UUID
);
-- +goose Down
DROP TABLE posts;

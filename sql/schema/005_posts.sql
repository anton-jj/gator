-- +goose Up
CREATE TABLE posts(
	id UUID NOT NULL PRIMARY KEY, 
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	title TEXT NOT NULL,
	url TEXT NOT NULL,
	description TEXT,
	published_at TIMESTAMPTZ,
	feed_id UUID NOT NULL
);
ALTER TABLE posts
ADD CONSTRAINT posts_url_key UNIQUE (url);
-- +goose Down
DROP TABLE posts;

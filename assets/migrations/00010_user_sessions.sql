-- +goose Up
CREATE TABLE user_sessions (
	id TEXT NOT NULL PRIMARY KEY,
	user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
	token TEXT NOT NULL,
	ip_address TEXT NOT NULL,
	user_agent TEXT NOT NULL,
	expires_at TIMESTAMP with time zone NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX user_sessions_token_idx ON user_sessions(token);

-- +goose Down
DROP TABLE user_sessions;
DROP INDEX user_sessions_token_idx;

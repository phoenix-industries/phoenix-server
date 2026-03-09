-- +goose Up
CREATE TYPE user_role AS ENUM ('root', 'admin', 'manager', 'member');
CREATE TABLE users (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	phone TEXT NOT NULL UNIQUE,
	role user_role NOT NULL DEFAULT 'member',
	city TEXT NOT NULL,
	governorate TEXT NOT NULL,
	address TEXT NOT NULL,
	birthdate TIMESTAMP WITH TIME ZONE NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX users_email_idx ON users (email);
CREATE INDEX users_phone_idx ON users (phone);

-- +goose Down
DROP INDEX users_email_idx;
DROP INDEX users_phone_idx;
DROP TABLE users;
DROP TYPE user_role;

-- +goose Up
ALTER TABLE users ADD COLUMN picture_id TEXT REFERENCES media(id);

-- +goose Down
ALTER TABLE users DROP COLUMN picture_id;

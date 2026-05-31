-- +goose Up
ALTER TABLE transactions ADD COLUMN status TEXT NOT NULL DEFAULT 'pending';

-- +goose Down
ALTER TABLE transactions DROP COLUMN status;

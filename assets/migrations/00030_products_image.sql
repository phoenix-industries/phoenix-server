-- +goose Up
ALTER TABLE products RENAME COLUMN media_id TO image_id;

-- +goose Down
ALTER TABLE products RENAME COLUMN image_id TO media_id;

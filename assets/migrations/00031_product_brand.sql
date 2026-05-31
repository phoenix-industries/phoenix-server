-- +goose Up
ALTER TABLE products ADD COLUMN brand TEXT;

-- +goose Down
ALTER TABLE products DROP COLUMN brand;

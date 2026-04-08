-- +goose Up
ALTER TABLE products ADD COLUMN price INTEGER;

-- +goose Down
ALTER TABLE products DROP COLUMN price;

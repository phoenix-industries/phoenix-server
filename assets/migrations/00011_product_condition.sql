-- +goose Up
ALTER TABLE products ADD COLUMN condition TEXT;

-- +goose Down
ALTER TABLE products DROP COLUMN condition;

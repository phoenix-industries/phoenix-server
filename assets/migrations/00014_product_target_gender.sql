-- +goose Up
ALTER TABLE products ADD COLUMN target_gender TEXT;

-- +goose Down
ALTER TABLE products DROP COLUMN target_gender;

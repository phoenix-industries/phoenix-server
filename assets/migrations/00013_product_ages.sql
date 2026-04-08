-- +goose Up
ALTER TABLE products ADD COLUMN minimum_age INTEGER;
ALTER TABLE products ADD COLUMN maximum_age INTEGER;

-- +goose Down
ALTER TABLE products DROP COLUMN minimum_age;
ALTER TABLE products DROP COLUMN maximum_age;

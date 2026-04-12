-- +goose Up
ALTER TABLE products ADD COLUMN description TEXT NOT null DEFAULT '';

-- +goose Down
ALTER TABLE products DROP COLUMN description;

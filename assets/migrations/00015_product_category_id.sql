-- +goose Up
ALTER TABLE products RENAME COLUMN category TO category_id;

-- +goose Down
ALTER TABLE products RENAME COLUMN category_id TO category;

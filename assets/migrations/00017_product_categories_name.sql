-- +goose Up
ALTER TABLE product_categories RENAME COLUMN category TO name;

-- +goose Down
ALTER TABLE product_categories RENAME COLUMN name TO category;

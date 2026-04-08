-- +goose Up
ALTER TABLE product_tags DROP CONSTRAINT product_tags_product_id_fkey;
DROP TABLE product_tags;
ALTER TABLE products ADD COLUMN tags TEXT;

-- +goose Down

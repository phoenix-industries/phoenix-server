-- +goose Up

ALTER TABLE products ADD COLUMN discount INTEGER NOT NULL DEFAULT 0;
ALTER TABLE products ADD COLUMN quantity INTEGER NOT NULL DEFAULT 0;

-- +goose Down

ALTER TABLE products DROP COLUMN discount;
ALTER TABLE products DROP COLUMN quantity;

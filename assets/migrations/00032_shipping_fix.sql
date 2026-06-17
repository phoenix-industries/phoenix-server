-- +goose Up
ALTER TABLE shippings ADD COLUMN full_name TEXT NOT NULL;
ALTER TABLE shippings ADD COLUMN phone TEXT NOT NULL;
ALTER TABLE shippings ADD COLUMN city TEXT NOT NULL;

-- +goose Down
ALTER TABLE shippings DROP COLUMN full_name;
ALTER TABLE shippings DROP COLUMN phone;
ALTER TABLE shippings DROP COLUMN city;

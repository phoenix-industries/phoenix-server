-- +goose Up
ALTER TABLE products ADD COLUMN reviewed_at TIMESTAMP WITH TIME ZONE DEFAULT null;
ALTER TABLE products ADD COLUMN approved BOOLEAN DEFAULT false;

-- +goose Down
ALTER TABLE products DROP COLUMN reviewed_at;
ALTER TABLE products DROP COLUMN approved;

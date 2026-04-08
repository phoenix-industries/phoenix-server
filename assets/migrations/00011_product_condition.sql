-- +goose Up
ALTER TABLE products ADD COLUMN condition TEXT;

-- +goose Down
SELECT 'down SQL query';
ALTER TABLE products DROP COLUMN condition;

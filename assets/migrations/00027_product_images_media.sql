-- +goose Up
ALTER TABLE product_images DROP COLUMN image;
ALTER TABLE product_images ADD COLUMN media_id TEXT REFERENCES media(id) ON DELETE CASCADE;
ALTER TABLE product_images DROP COLUMN product_id;
ALTER TABLE product_images ADD COLUMN product_id TEXT REFERENCES products(id) ON DELETE CASCADE;

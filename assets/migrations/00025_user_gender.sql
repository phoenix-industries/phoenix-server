-- +goose Up

ALTER TABLE users ADD COLUMN gender TEXT NOT NULL;

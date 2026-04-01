-- +goose Up
ALTER TABLE users ADD password TEXT NOT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN password;

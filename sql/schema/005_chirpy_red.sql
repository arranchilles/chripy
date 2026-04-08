-- +goose Up
ALTER TABLE users ADD is_chirp_red BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE users DROP COLUMN is_chirp_red;

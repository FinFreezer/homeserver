-- +goose Up
ALTER TABLE users
ADD COLUMN authtoken TEXT UNIQUE NOT NULL DEFAULT 'unset';

-- +goose Down
ALTER TABLE users
DROP COLUMN authtoken;

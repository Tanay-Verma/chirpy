-- +goose Up
CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  email TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS users;

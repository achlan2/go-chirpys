-- +goose Up
CREATE TABLE users (
  id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  email TEXT NOT NULL UNIQUE,
  hashed_password TEXT NOT NULL
);
-- +goose Down
DROP TABLE users;

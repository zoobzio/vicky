-- +goose Up
CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL,
    name TEXT,
    avatar_url TEXT,
    access_token TEXT NOT NULL,
    token_expiry TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_login_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE users;

-- +goose Up
CREATE TABLE keys (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    key_hash   TEXT NOT NULL UNIQUE,
    key_prefix TEXT NOT NULL,
    scopes     TEXT[] NOT NULL DEFAULT '{}',
    rate_limit INTEGER,
    expires_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_keys_user_id ON keys(user_id);
CREATE INDEX idx_keys_key_hash ON keys(key_hash);

-- +goose Down
DROP TABLE keys;

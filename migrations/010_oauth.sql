-- +goose Up
-- OAuth state for CSRF protection
CREATE TABLE oauth_states (
    id BIGSERIAL PRIMARY KEY,
    state TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_oauth_states_state ON oauth_states(state);
CREATE INDEX idx_oauth_states_expires_at ON oauth_states(expires_at);

-- Sessions for cookie-based auth
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    tenant_id TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL DEFAULT '',
    scopes TEXT[],
    roles TEXT[],
    meta JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- +goose Down
DROP TABLE sessions;
DROP TABLE oauth_states;

-- +goose Up
CREATE TABLE repositories (
    id BIGSERIAL PRIMARY KEY,
    github_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    owner TEXT NOT NULL,
    name TEXT NOT NULL,
    full_name TEXT NOT NULL,
    description TEXT,
    default_branch TEXT NOT NULL DEFAULT 'main',
    private BOOLEAN NOT NULL DEFAULT false,
    html_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, github_id)
);

CREATE INDEX idx_repositories_user_id ON repositories(user_id);

-- +goose Down
DROP TABLE repositories;

-- +goose Up
CREATE TABLE versions (
    id BIGSERIAL PRIMARY KEY,
    repository_id BIGINT NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    owner TEXT NOT NULL,
    repo_name TEXT NOT NULL,
    tag TEXT NOT NULL,
    commit_sha TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (repository_id, tag),
    UNIQUE (user_id, owner, repo_name, tag)
);

CREATE INDEX idx_versions_repository_id ON versions(repository_id);
CREATE INDEX idx_versions_user_id ON versions(user_id);
CREATE INDEX idx_versions_status ON versions(status);
CREATE INDEX idx_versions_lookup ON versions(user_id, owner, repo_name);

-- +goose Down
DROP TABLE versions;

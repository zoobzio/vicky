-- +goose Up
CREATE TABLE jobs (
    id BIGSERIAL PRIMARY KEY,
    version_id BIGINT NOT NULL REFERENCES versions(id) ON DELETE CASCADE,
    repository_id BIGINT NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    owner TEXT NOT NULL,
    repo_name TEXT NOT NULL,
    tag TEXT NOT NULL,
    stage TEXT NOT NULL DEFAULT 'fetch',
    status TEXT NOT NULL DEFAULT 'pending',
    progress INT NOT NULL DEFAULT 0,
    error TEXT,
    items_total INT NOT NULL DEFAULT 0,
    items_processed INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_jobs_version_id ON jobs(version_id);
CREATE INDEX idx_jobs_repository_id ON jobs(repository_id);
CREATE INDEX idx_jobs_user_id ON jobs(user_id);
CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_lookup ON jobs(user_id, owner, repo_name, tag);

-- +goose Down
DROP TABLE jobs;

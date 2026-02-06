-- +goose Up
CREATE TABLE documents (
    id BIGSERIAL PRIMARY KEY,
    version_id BIGINT NOT NULL REFERENCES versions(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    owner TEXT NOT NULL,
    repo_name TEXT NOT NULL,
    tag TEXT NOT NULL,
    path TEXT NOT NULL,
    content_type TEXT NOT NULL,
    content_hash TEXT NOT NULL,
    vector vector(1536),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (version_id, path),
    UNIQUE (user_id, owner, repo_name, tag, path)
);

CREATE INDEX idx_documents_version_id ON documents(version_id);
CREATE INDEX idx_documents_user_id ON documents(user_id);
CREATE INDEX idx_documents_content_type ON documents(content_type);
CREATE INDEX idx_documents_lookup ON documents(user_id, owner, repo_name, tag);

-- +goose Down
DROP TABLE documents;

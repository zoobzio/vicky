-- +goose Up
CREATE TABLE chunks (
    id BIGSERIAL PRIMARY KEY,
    document_id BIGINT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    owner TEXT NOT NULL,
    repo_name TEXT NOT NULL,
    tag TEXT NOT NULL,
    path TEXT NOT NULL,
    kind TEXT NOT NULL,
    start_line INT NOT NULL,
    end_line INT NOT NULL,
    symbol TEXT,
    context TEXT[],
    content TEXT NOT NULL,
    vector vector(1536),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_chunks_document_id ON chunks(document_id);
CREATE INDEX idx_chunks_user_id ON chunks(user_id);
CREATE INDEX idx_chunks_kind ON chunks(kind);
CREATE INDEX idx_chunks_lookup ON chunks(user_id, owner, repo_name, tag);
CREATE INDEX idx_chunks_path_lookup ON chunks(user_id, owner, repo_name, tag, path);

-- +goose Down
DROP TABLE chunks;

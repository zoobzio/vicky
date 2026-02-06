-- +goose Up
CREATE TABLE symbols (
    id BIGSERIAL PRIMARY KEY,
    version_id BIGINT NOT NULL REFERENCES versions(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    owner TEXT NOT NULL,
    repo_name TEXT NOT NULL,
    tag TEXT NOT NULL,
    name TEXT NOT NULL,
    qualified_name TEXT NOT NULL,
    kind TEXT NOT NULL,
    signature TEXT,
    doc TEXT,
    file_path TEXT NOT NULL,
    start_line INT NOT NULL,
    end_line INT NOT NULL,
    exported BOOLEAN NOT NULL DEFAULT false,
    parent_id BIGINT REFERENCES symbols(id) ON DELETE SET NULL,
    vector vector(1536),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_symbols_version_id ON symbols(version_id);
CREATE INDEX idx_symbols_user_id ON symbols(user_id);
CREATE INDEX idx_symbols_kind ON symbols(kind);
CREATE INDEX idx_symbols_exported ON symbols(exported);
CREATE INDEX idx_symbols_parent_id ON symbols(parent_id);
CREATE INDEX idx_symbols_lookup ON symbols(user_id, owner, repo_name, tag);

-- +goose Down
DROP TABLE symbols;

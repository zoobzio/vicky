-- +goose Up
CREATE TABLE ingestion_configs (
    id BIGSERIAL PRIMARY KEY,
    repository_id BIGINT NOT NULL UNIQUE REFERENCES repositories(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    language TEXT NOT NULL CHECK (language IN ('go', 'typescript')),
    include_docs BOOLEAN NOT NULL DEFAULT true,
    exclude_patterns TEXT[] DEFAULT '{}',
    max_file_size BIGINT NOT NULL DEFAULT 1048576,
    language_config JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_ingestion_configs_repository_id ON ingestion_configs(repository_id);
CREATE INDEX idx_ingestion_configs_user_id ON ingestion_configs(user_id);
CREATE INDEX idx_ingestion_configs_language ON ingestion_configs(language);

-- +goose Down
DROP TABLE ingestion_configs;

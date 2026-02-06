-- +goose Up
-- Update capacitor configs with proper defaults for pipeline stages.
-- Timeouts stored as nanoseconds (Go time.Duration is int64 nanoseconds).

-- Fetch stage defaults: 8 workers, 30s timeout
INSERT INTO configs (domain, data) VALUES
    ('fetch', '{"workers": 8, "timeout": 30000000000}')
ON CONFLICT (domain) DO UPDATE SET
    data = EXCLUDED.data,
    updated_at = now();

-- Parse stage defaults: 8 workers, 60s timeout
INSERT INTO configs (domain, data) VALUES
    ('parse', '{"workers": 8, "timeout": 60000000000}')
ON CONFLICT (domain) DO UPDATE SET
    data = EXCLUDED.data,
    updated_at = now();

-- Chunk stage defaults: 8 workers, 30s timeout
INSERT INTO configs (domain, data) VALUES
    ('chunk', '{"workers": 8, "timeout": 30000000000}')
ON CONFLICT (domain) DO UPDATE SET
    data = EXCLUDED.data,
    updated_at = now();

-- Embedding stage defaults: 4 workers, 128 batch size, 30s timeout
INSERT INTO configs (domain, data) VALUES
    ('embedding', '{"workers": 4, "batch_size": 128, "timeout": 30000000000}')
ON CONFLICT (domain) DO UPDATE SET
    data = EXCLUDED.data,
    updated_at = now();

-- Observability defaults (empty - configured via database)
INSERT INTO configs (domain, data) VALUES
    ('observability', '{}')
ON CONFLICT (domain) DO NOTHING;

-- +goose Down
DELETE FROM configs WHERE domain IN ('fetch', 'parse', 'chunk');
-- Restore old embedding schema
UPDATE configs SET data = '{"code": {"name": "stub", "dimensions": 1536}, "docs": {"name": "stub", "dimensions": 1536}}' WHERE domain = 'embedding';

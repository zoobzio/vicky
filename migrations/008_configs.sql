-- +goose Up
-- Runtime configuration storage for hot-reload via flux capacitors.
-- Uses LISTEN/NOTIFY for real-time change propagation.

CREATE TABLE configs (
    domain TEXT PRIMARY KEY,
    data JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION notify_config_change() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('config_' || NEW.domain, NEW.domain);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER config_change_trigger
    AFTER INSERT OR UPDATE ON configs
    FOR EACH ROW
    EXECUTE FUNCTION notify_config_change();

-- Seed default configs
INSERT INTO configs (domain, data) VALUES
    ('embedding', '{"code": {"name": "stub", "dimensions": 1536}, "docs": {"name": "stub", "dimensions": 1536}}'),
    ('events', '{"signals": {}}')
ON CONFLICT (domain) DO NOTHING;

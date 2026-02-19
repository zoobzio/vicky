-- +goose Up
-- Recreate check constraint to include new cancellation statuses
ALTER TABLE jobs DROP CONSTRAINT IF EXISTS jobs_status_check;
ALTER TABLE jobs ADD CONSTRAINT jobs_status_check
    CHECK (status IN ('pending', 'running', 'completed', 'failed', 'cancelling', 'cancelled'));

-- +goose Down
-- Restore original check constraint
ALTER TABLE jobs DROP CONSTRAINT IF EXISTS jobs_status_check;
ALTER TABLE jobs ADD CONSTRAINT jobs_status_check
    CHECK (status IN ('pending', 'running', 'completed', 'failed'));

-- Drop indexes
DROP INDEX IF EXISTS idx_jobs_run_id;
DROP INDEX IF EXISTS idx_workflow_runs_status;
DROP INDEX IF EXISTS idx_workflow_runs_workflow_id;

-- Revert job columns
ALTER TABLE jobs DROP COLUMN IF EXISTS run_id;
ALTER TABLE jobs DROP COLUMN IF EXISTS lease_expires_at;
ALTER TABLE jobs DROP COLUMN IF EXISTS result;
ALTER TABLE jobs DROP COLUMN IF EXISTS error;

-- Drop workflow_runs table and enum
DROP TABLE IF EXISTS workflow_runs;
DROP TYPE IF EXISTS run_status;

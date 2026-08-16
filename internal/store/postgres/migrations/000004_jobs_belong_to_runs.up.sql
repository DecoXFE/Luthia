-- Jobs belong to runs. A job's workflow is reached through its run,
-- so the redundant workflow_id column on jobs is removed.

DROP INDEX IF EXISTS idx_jobs_workflow_id;

ALTER TABLE jobs DROP COLUMN IF EXISTS workflow_id;
ALTER TABLE jobs ALTER COLUMN run_id SET NOT NULL;

-- Restore workflow_id on jobs, backfilling it from the owning run
ALTER TABLE jobs ALTER COLUMN run_id DROP NOT NULL;

ALTER TABLE jobs ADD COLUMN workflow_id UUID REFERENCES workflows(id) ON DELETE CASCADE;
UPDATE jobs
SET workflow_id = wr.workflow_id
FROM workflow_runs wr
WHERE jobs.run_id = wr.id;

ALTER TABLE jobs ALTER COLUMN workflow_id SET NOT NULL;
CREATE INDEX idx_jobs_workflow_id ON jobs(workflow_id);

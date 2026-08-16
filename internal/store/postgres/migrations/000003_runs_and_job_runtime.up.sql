-- Create run status enum
CREATE TYPE run_status AS ENUM ('PENDING', 'RUNNING', 'COMPLETED', 'FAILED', 'CANCELLED');

-- Create workflow_runs table: one execution of a workflow
CREATE TABLE workflow_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    status run_status NOT NULL DEFAULT 'PENDING',
    error TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    finished_at TIMESTAMP WITH TIME ZONE
);

-- Add runtime fields to jobs
ALTER TABLE jobs ADD COLUMN run_id UUID REFERENCES workflow_runs(id) ON DELETE CASCADE;
ALTER TABLE jobs ADD COLUMN lease_expires_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE jobs ADD COLUMN result JSONB;
ALTER TABLE jobs ADD COLUMN error TEXT;

-- Indexes for performance
CREATE INDEX idx_workflow_runs_workflow_id ON workflow_runs(workflow_id);
CREATE INDEX idx_workflow_runs_status ON workflow_runs(status);
CREATE INDEX idx_jobs_run_id ON jobs(run_id);

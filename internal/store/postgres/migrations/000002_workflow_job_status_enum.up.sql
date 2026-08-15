-- Create status enums
CREATE TYPE workflow_status AS ENUM ('active', 'inactive');
CREATE TYPE job_status AS ENUM ('CREATED', 'QUEUED', 'PICKED_UP', 'RUNNING', 'COMPLETED', 'FAILED', 'RETRYING', 'DEAD_LETTER', 'CANCELLED');

-- Drop defaults before altering column types
ALTER TABLE workflows ALTER COLUMN status DROP DEFAULT;
ALTER TABLE jobs ALTER COLUMN status DROP DEFAULT;

-- Convert existing status columns to enum types
ALTER TABLE workflows ALTER COLUMN status TYPE workflow_status USING status::workflow_status;
ALTER TABLE jobs ALTER COLUMN status TYPE job_status USING status::job_status;

-- Restore defaults with enum values
ALTER TABLE workflows ALTER COLUMN status SET DEFAULT 'active'::workflow_status;
ALTER TABLE jobs ALTER COLUMN status SET DEFAULT 'CREATED'::job_status;

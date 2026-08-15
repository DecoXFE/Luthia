-- Revert status columns back to varchar
ALTER TABLE jobs ALTER COLUMN status TYPE VARCHAR(50) USING status::text;
ALTER TABLE workflows ALTER COLUMN status TYPE VARCHAR(50) USING status::text;

-- Drop enum types
DROP TYPE IF EXISTS job_status;
DROP TYPE IF EXISTS workflow_status;

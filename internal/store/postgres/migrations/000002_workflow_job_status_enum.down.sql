-- Revert status columns back to varchar
-- Drop defaults first so the enum type has no dependents
ALTER TABLE jobs ALTER COLUMN status DROP DEFAULT;
ALTER TABLE workflows ALTER COLUMN status DROP DEFAULT;

-- Convert columns to varchar
ALTER TABLE jobs ALTER COLUMN status TYPE VARCHAR(50) USING status::text;
ALTER TABLE workflows ALTER COLUMN status TYPE VARCHAR(50) USING status::text;

-- Restore plain-text defaults
ALTER TABLE jobs ALTER COLUMN status SET DEFAULT 'CREATED';
ALTER TABLE workflows ALTER COLUMN status SET DEFAULT 'active';

-- Drop enum types
DROP TYPE IF EXISTS job_status;
DROP TYPE IF EXISTS workflow_status;

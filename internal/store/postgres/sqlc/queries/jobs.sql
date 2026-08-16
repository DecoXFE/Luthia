-- name: CreateJob :one
INSERT INTO jobs (run_id, step, payload)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetJob :one
SELECT *
FROM jobs
WHERE id = $1;

-- name: ListJobsByRun :many
SELECT *
FROM jobs
WHERE run_id = $1
ORDER BY created_at;

-- name: ClaimNextJob :one
WITH claimed AS (
    SELECT id
    FROM jobs
    WHERE status = 'QUEUED'
      AND (scheduled_at IS NULL OR scheduled_at <= NOW())
    ORDER BY priority DESC NULLS LAST, created_at ASC
    FOR UPDATE SKIP LOCKED
    LIMIT 1
)
UPDATE jobs
SET status = 'RUNNING',
    lease_expires_at = $1,
    updated_at = NOW()
WHERE id IN (SELECT id FROM claimed)
RETURNING *;

-- name: CompleteJob :one
UPDATE jobs
SET status = 'COMPLETED',
    result = $2,
    error = NULL,
    lease_expires_at = NULL,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: FailJob :one
UPDATE jobs
SET status = $2,
    error = $3,
    attempt_count = attempt_count + 1,
    lease_expires_at = NULL,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: HeartbeatJob :one
UPDATE jobs
SET lease_expires_at = $2
WHERE id = $1
  AND status = 'RUNNING'
RETURNING *;

-- name: RequeueJob :one
UPDATE jobs
SET status = 'QUEUED',
    lease_expires_at = NULL,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: PromoteRetries :many
UPDATE jobs
SET status = 'QUEUED',
    lease_expires_at = NULL,
    updated_at = NOW()
WHERE status = 'RETRYING'
  AND scheduled_at IS NOT NULL
  AND scheduled_at <= NOW()
RETURNING *;

-- name: FindExpiredLeases :many
SELECT *
FROM jobs
WHERE status = 'RUNNING'
  AND lease_expires_at IS NOT NULL
  AND lease_expires_at < NOW();

-- name: CreateJobEvent :one
INSERT INTO job_events (job_id, event_type, data)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListJobEvents :many
SELECT *
FROM job_events
WHERE job_id = $1
ORDER BY created_at;

-- name: CreateRun :one
INSERT INTO workflow_runs (workflow_id, status)
VALUES ($1, 'RUNNING')
RETURNING *;

-- name: GetRun :one
SELECT *
FROM workflow_runs
WHERE id = $1;

-- name: ListRunsByWorkflow :many
SELECT *
FROM workflow_runs
WHERE workflow_id = $1
ORDER BY created_at DESC;

-- name: CompleteRun :one
UPDATE workflow_runs
SET status = 'COMPLETED',
    finished_at = NOW()
WHERE id = $1
RETURNING *;

-- name: FailRun :one
UPDATE workflow_runs
SET status = 'FAILED',
    error = $2,
    finished_at = NOW()
WHERE id = $1
RETURNING *;

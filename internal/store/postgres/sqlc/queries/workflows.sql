-- name: ListWorkflows :many
SELECT *
FROM workflows;

-- name: CreateWorkflow :one
INSERT INTO workflows (name, description)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteWorkflow :one
DELETE FROM workflows
WHERE id = $1
RETURNING id;

-- name: GetWorkflow :one
SELECT *
FROM workflows
WHERE id = $1;
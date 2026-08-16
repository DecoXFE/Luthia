# Jobs

:::note[Planned]

The Jobs API is **not implemented yet**. The `jobs` table, the `job_status` enum (`CREATED`, `QUEUED`, `PICKED_UP`, `RUNNING`, `COMPLETED`, `FAILED`, `RETRYING`, `DEAD_LETTER`, `CANCELLED`) and the `job_events` table exist in the schema, but there are no endpoints or worker to process them yet. This page documents the target behavior.

:::

A **job** is a single unit of work within a workflow. This is what actually gets executed.

## Submitting a Job

```bash
curl -X POST http://localhost:8080/api/workflows/<workflow-id>/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "step": "resize",
    "payload": {
      "image": "photo1.jpg",
      "width": 800,
      "height": 600
    }
  }'
```

Response:

```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "workflow_id": "550e8400-e29b-41d4-a716-446655440000",
  "step": "resize",
  "payload": {"image": "photo1.jpg", "width": 800, "height": 600},
  "status": "QUEUED",
  "attempt_count": 0,
  "max_retries": 3,
  "created_at": "2025-01-15T10:35:00Z"
}
```

## Job Structure

| Field | Type | Description |
|-------|------|-------------|
| `id` | UUID | Unique identifier |
| `workflow_id` | UUID | Parent workflow |
| `step` | string | What kind of work (e.g., "resize", "send-email") |
| `payload` | JSONB | Data needed to execute |
| `status` | string | Current state |
| `priority` | int | Higher = processed first |
| `attempt_count` | int | How many times we tried |
| `max_retries` | int | Maximum attempts before dead letter |
| `scheduled_at` | timestamp | Delayed execution (optional) |

## Job Statuses

```
CREATED → QUEUED → PICKED_UP → RUNNING → COMPLETED
                                ↓
                            FAILED → RETRYING → QUEUED
                                ↓
                            DEAD_LETTER
                                ↓
                            CANCELLED
```

### Status Meaning

| Status | Meaning |
|--------|---------|
| `CREATED` | Job exists but not yet queued |
| `QUEUED` | In Redis, waiting for a worker |
| `PICKED_UP` | Worker claimed it |
| `RUNNING` | Worker is executing |
| `COMPLETED` | Done successfully |
| `FAILED` | Execution failed |
| `RETRYING` | Will be retried |
| `DEAD_LETTER` | Max retries exceeded |
| `CANCELLED` | Manually cancelled |

## Retry Logic

When a job fails:

```
Attempt 1: FAILED → RETRYING → QUEUED (wait 2s)
Attempt 2: FAILED → RETRYING → QUEUED (wait 4s)
Attempt 3: FAILED → DEAD_LETTER
```

The wait time increases exponentially (backoff).

## Manual Retry

If a dead letter job should be retried:

```bash
curl -X POST http://localhost:8080/api/jobs/<job-id>/retry
```

## Cancel a Job

```bash
curl -X POST http://localhost:8080/api/jobs/<job-id>/cancel
```

## Job Events

Every state change is logged:

```bash
curl http://localhost:8080/api/jobs/<job-id>/events
```

Response:

```json
[
  {"event_type": "CREATED", "created_at": "2025-01-15T10:35:00Z"},
  {"event_type": "QUEUED", "created_at": "2025-01-15T10:35:01Z"},
  {"event_type": "PICKED_UP", "data": {"worker": "worker-1"}, "created_at": "2025-01-15T10:35:02Z"},
  {"event_type": "RUNNING", "created_at": "2025-01-15T10:35:02Z"},
  {"event_type": "COMPLETED", "data": {"duration_ms": 1250}, "created_at": "2025-01-15T10:35:03Z"}
]
```

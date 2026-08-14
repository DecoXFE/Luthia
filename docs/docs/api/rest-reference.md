# REST API Reference

Base URL: `http://localhost:8080`

## Workflows

### Create Workflow

```
POST /api/workflows
```

**Request:**

```json
{
  "name": "process-images",
  "description": "Resize and compress images",
  "config": {}
}
```

**Response:** `201 Created`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "process-images",
  "description": "Resize and compress images",
  "status": "active",
  "config": {},
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

### List Workflows

```
GET /api/workflows
```

**Response:** `200 OK`

```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "process-images",
    "status": "active",
    "created_at": "2025-01-15T10:30:00Z"
  }
]
```

### Get Workflow

```
GET /api/workflows/:id
```

**Response:** `200 OK`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "process-images",
  "description": "Resize and compress images",
  "status": "active",
  "config": {},
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

---

## Jobs

### Submit Job

```
POST /api/workflows/:workflow_id/jobs
```

**Request:**

```json
{
  "step": "resize",
  "payload": {
    "image": "photo.jpg",
    "width": 800
  },
  "priority": 0,
  "max_retries": 3
}
```

**Response:** `201 Created`

```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "workflow_id": "550e8400-e29b-41d4-a716-446655440000",
  "step": "resize",
  "payload": {"image": "photo.jpg", "width": 800},
  "status": "QUEUED",
  "attempt_count": 0,
  "max_retries": 3,
  "created_at": "2025-01-15T10:35:00Z"
}
```

### List Jobs

```
GET /api/workflows/:workflow_id/jobs
```

**Query Parameters:**

| Param | Description |
|-------|-------------|
| `status` | Filter by status (e.g., `RUNNING`, `FAILED`) |
| `limit` | Max results (default: 100) |
| `offset` | Pagination offset |

### Get Job

```
GET /api/jobs/:id
```

### Get Job Events

```
GET /api/jobs/:id/events
```

**Response:** `200 OK`

```json
[
  {
    "id": "event-1",
    "job_id": "660e8400-e29b-41d4-a716-446655440001",
    "event_type": "CREATED",
    "data": {},
    "created_at": "2025-01-15T10:35:00Z"
  },
  {
    "id": "event-2",
    "job_id": "660e8400-e29b-41d4-a716-446655440001",
    "event_type": "QUEUED",
    "data": {},
    "created_at": "2025-01-15T10:35:01Z"
  }
]
```

### Retry Job

```
POST /api/jobs/:id/retry
```

Re-queues a failed or dead letter job.

### Cancel Job

```
POST /api/jobs/:id/cancel
```

Cancels a running or queued job.

---

## Workers

### List Workers

```
GET /api/workers
```

**Response:** `200 OK`

```json
[
  {
    "id": "worker-abc",
    "hostname": "server-1",
    "status": "busy",
    "current_job_id": "660e8400-e29b-41d4-a716-446655440001",
    "last_heartbeat": "2025-01-15T10:40:00Z"
  }
]
```

---

## Stats

### Get Dashboard Stats

```
GET /api/stats
```

**Response:** `200 OK`

```json
{
  "workflows": 5,
  "jobs": {
    "total": 10247,
    "running": 42,
    "completed": 10150,
    "failed": 55
  },
  "workers": {
    "active": 3,
    "idle": 1,
    "offline": 0
  }
}
```

---

## Health Check

```
GET /health
```

**Response:** `200 OK`

```json
{
  "status": "ok"
}
```

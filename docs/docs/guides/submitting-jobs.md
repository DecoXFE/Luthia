# Submitting Jobs

Learn how to submit jobs to Luthia for processing.

## Submit a Single Job

```bash
curl -X POST http://localhost:8080/api/workflows/<workflow-id>/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "step": "resize",
    "payload": {
      "image": "photo.jpg",
      "width": 800
    }
  }'
```

## Submit Multiple Jobs

```bash
for img in photo1.jpg photo2.jpg photo3.jpg; do
  curl -X POST http://localhost:8080/api/workflows/<workflow-id>/jobs \
    -H "Content-Type: application/json" \
    -d "{\"step\": \"resize\", \"payload\": {\"image\": \"$img\"}}"
done
```

## Job Options

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `step` | string | required | Task type to execute |
| `payload` | object | `{}` | Data for the job |
| `priority` | int | `0` | Higher = processed first |
| `max_retries` | int | `3` | Max attempts before dead letter |
| `scheduled_at` | timestamp | now | Delayed execution |

## Delayed Jobs

Execute a job after a delay:

```bash
curl -X POST http://localhost:8080/api/workflows/<workflow-id>/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "step": "send-email",
    "payload": {"to": "user@example.com"},
    "scheduled_at": "2025-01-15T12:00:00Z"
  }'
```

## Priority Jobs

Higher priority jobs are processed first:

```bash
# Normal priority (0)
curl -X POST http://localhost:8080/api/workflows/<workflow-id>/jobs \
  -d '{"step": "process", "payload": {"data": "normal"}}'

# High priority (10)
curl -X POST http://localhost:8080/api/workflows/<workflow-id>/jobs \
  -d '{"step": "process", "payload": {"data": "urgent"}, "priority": 10}'
```

## Check Job Status

```bash
curl http://localhost:8080/api/jobs/<job-id>
```

## View Job History

```bash
curl http://localhost:8080/api/jobs/<job-id>/events
```

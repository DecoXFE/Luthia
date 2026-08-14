# Configuring Retries

Luthia automatically retries failed jobs. Here's how to configure this behavior.

## Default Behavior

- **Max retries**: 3
- **Backoff**: Exponential (2s, 4s, 8s, ...)

## Custom Retry Per Job

```bash
curl -X POST http://localhost:8080/api/workflows/<workflow-id>/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "step": "send-email",
    "payload": {"to": "user@example.com"},
    "max_retries": 5
  }'
```

## No Retries

```bash
curl -X POST http://localhost:8080/api/workflows/<workflow-id>/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "step": "critical-operation",
    "payload": {"data": "important"},
    "max_retries": 0
  }'
```

## Retry Flow

```
Attempt 1: FAILED
    ↓ (wait 2s)
Attempt 2: FAILED
    ↓ (wait 4s)
Attempt 3: FAILED
    ↓
DEAD_LETTER (no more retries)
```

## Manual Retry

For dead letter jobs, you can manually retry:

```bash
curl -X POST http://localhost:8080/api/jobs/<job-id>/retry
```

This resets the attempt count and re-queues the job.

## Dead Letter Queue

When a job exceeds max retries, it enters the dead letter queue. This means:

- It won't be automatically retried
- It's preserved for inspection
- You can manually retry or cancel it

## Best Practices

1. **Use retries for transient failures** — Network timeouts, temporary overloads
2. **Don't retry for permanent failures** — Invalid data, permission errors
3. **Set lower max retries for critical jobs** — Avoid wasting resources
4. **Monitor dead letters** — They indicate real problems
5. **Use manual retry after fixing the root cause**

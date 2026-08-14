# Creating Workflows

This guide shows you how to create and manage workflows in Luthia.

## Basic Workflow

```bash
curl -X POST http://localhost:8080/api/workflows \
  -H "Content-Type: application/json" \
  -d '{"name": "my-workflow", "description": "Does something useful"}'
```

## Workflow with Configuration

```bash
curl -X POST http://localhost:8080/api/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "image-processing",
    "description": "Resize and compress images",
    "config": {
      "default_max_retries": 5,
      "timeout_seconds": 300
    }
  }'
```

## List Your Workflows

```bash
curl http://localhost:8080/api/workflows
```

## Get a Specific Workflow

```bash
curl http://localhost:8080/api/workflows/<workflow-id>
```

## Pausing a Workflow

```bash
curl -X PATCH http://localhost:8080/api/workflows/<workflow-id> \
  -H "Content-Type: application/json" \
  -d '{"status": "paused"}'
```

When paused, new jobs can still be submitted but workers won't pick them up until the workflow is active again.

## Best Practices

1. **One workflow per job type** — Don't mix unrelated jobs
2. **Descriptive names** — `process-images` not `workflow-1`
3. **Use config for defaults** — Set default retries, timeouts
4. **Monitor via dashboard** — Check progress regularly

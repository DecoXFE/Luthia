# Creating Workflows

This guide shows you how to create and manage workflows in Luthia.

## Basic Workflow

```bash
curl -X POST http://localhost:8080/api/workflows \
  -H "Content-Type: application/json" \
  -d '{"name": "my-workflow", "description": "Does something useful"}'
```

The only required field is `name`. `description` is optional.

## List Your Workflows

```bash
curl http://localhost:8080/api/workflows
```

## Delete a Workflow

```bash
curl -X DELETE http://localhost:8080/api/workflows/<workflow-id>
```

Returns `204 No Content` on success. This permanently deletes the workflow and its jobs.

## Error Handling

| Status | Meaning |
|--------|---------|
| `400` | Invalid JSON body, or `name` missing / too long |
| `409` | A workflow with that `name` already exists |
| `404` | Workflow id doesn't exist (DELETE) |

## Workflow Configuration

:::note[Planned]

Workflow `config` (JSONB) is stored and returned, but setting it via the API is not implemented yet. It currently defaults to `{}`.

:::

## Best Practices

1. **One workflow per job type** — Don't mix unrelated jobs
2. **Descriptive names** — `process-images` not `workflow-1`
3. **Name uniqueness** — `name` is unique; use it as your external reference to the workflow
4. **Monitor via dashboard** — Check progress regularly

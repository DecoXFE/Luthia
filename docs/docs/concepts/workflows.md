# Workflows

A **workflow** is a named container for related jobs. It defines what kind of work you want to do.

## Creating a Workflow

```bash
curl -X POST http://localhost:8080/api/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "process-images",
    "description": "Resize and compress product images"
  }'
```

Response:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "process-images",
  "description": "Resize and compress product images",
  "status": "active",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

## Listing Workflows

```bash
curl http://localhost:8080/api/workflows
```

## Workflow Structure

| Field | Type | Description |
|-------|------|-------------|
| `id` | UUID | Unique identifier |
| `name` | string | Human-readable name (unique) |
| `description` | string | What this workflow does |
| `status` | string | `active` or `paused` |
| `config` | JSONB | Workflow-specific configuration |
| `created_at` | timestamp | When it was created |
| `updated_at` | timestamp | Last modification |

## Workflow = Organization

Think of workflows as folders for your jobs:

```
process-images workflow
├── Job: resize photo1.jpg
├── Job: resize photo2.jpg
├── Job: compress photo1.jpg
└── Job: compress photo2.jpg

send-emails workflow
├── Job: send to alice@example.com
├── Job: send to bob@example.com
└── Job: send to carol@example.com
```

Each workflow tracks its own jobs independently.

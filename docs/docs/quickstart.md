# Quickstart

Get Luthia running on your machine in 5 minutes.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- [Go](https://go.dev/dl/) 1.22+ (for local development)

## 1. Start the Database

Luthia currently uses Postgres as its store. Start just the database:

```bash
git clone https://github.com/DecoXFE/luthia.git
cd luthia
docker compose up postgres
```

## 2. Run Migrations

```bash
make migrate
```

This applies the SQL migrations (schema + enums) to the database.

## 3. Start the API

```bash
make run-api
```

The API server starts on `http://localhost:8080`.

## 4. Verify It Works

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{"status": "ok"}
```

## 5. Create Your First Workflow

```bash
curl -X POST http://localhost:8080/api/workflows \
  -H "Content-Type: application/json" \
  -d '{"name": "process-images", "description": "Resize and compress images"}'
```

You should get a `201 Created` response with the workflow's `id`.

## 6. List Workflows

```bash
curl http://localhost:8080/api/workflows
```

## Stop Everything

```bash
docker compose down
```

## What's Implemented Today

```
1. docker compose up postgres  → Postgres running
2. make migrate                → Schema + enums applied
3. make run-api                → API on :8080
4. POST /api/workflows         → Creates a workflow (201)
5. GET  /api/workflows         → Lists workflows
6. DELETE /api/workflows/:id   → Deletes a workflow (204)
```

:::note[Roadmap]

Jobs, workers, Redis queue and the dashboard are **planned** but not implemented yet.

:::

## Next Steps

- [Architecture](architecture) — Understand how the pieces fit together
- [Concepts: Workflows](concepts/workflows) — What workflows are and how to use them
- [API Reference](/api-reference) — All available endpoints (interactive)

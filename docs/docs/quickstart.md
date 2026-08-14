# Quickstart

Get Luthia running on your machine in 5 minutes.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- [Go](https://go.dev/dl/) 1.22+ (for local development)

## 1. Start Everything

```bash
git clone https://github.com/DecoXFE/luthia.git
cd luthia
make dev
```

This starts:
- **API server** on `http://localhost:8080`
- **PostgreSQL** on `localhost:5432`
- **Redis** on `localhost:6379`

## 2. Verify It Works

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{"status": "ok"}
```

## 3. Create Your First Workflow

```bash
curl -X POST http://localhost:8080/api/workflows \
  -H "Content-Type: application/json" \
  -d '{"name": "process-images", "description": "Resize and compress images"}'
```

## 4. Submit Jobs

```bash
curl -X POST http://localhost:8080/api/workflows/<workflow-id>/jobs \
  -H "Content-Type: application/json" \
  -d '{"step": "resize", "payload": {"image": "photo1.jpg", "width": 800}}'
```

## 5. Check Job Status

```bash
curl http://localhost:8080/api/jobs/<job-id>
```

## Stop Everything

```bash
make dev-down
```

## What Just Happened

```
1. make dev           → Docker builds Go binary, starts Postgres + Redis
2. POST /workflows    → Creates a workflow in the database
3. POST /jobs         → Creates a job with status QUEUED
4. Worker picks up    → Changes status to RUNNING, executes, marks COMPLETED
5. GET /jobs          → Shows current status and attempt history
```

## Next Steps

- [Architecture](architecture) — Understand how the pieces fit together
- [Concepts: Workflows](concepts/workflows) — What workflows are and how to use them
- [API Reference](api/rest-reference) — All available endpoints

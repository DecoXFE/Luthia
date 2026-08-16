# Architecture

This page explains how Luthia works internally — the components, how they communicate, and why each decision was made.

:::note[Status]

Today Luthia implements the **API server + Postgres store** for workflows (create, list, delete). The Redis queue, workers, jobs API and dashboard are **planned**; they appear below as the target architecture.

:::

## High-Level Overview

```
┌─────────────────────────────────────────────────┐
│                    CLIENTS                       │
│         (Your app, SDK, curl, dashboard)         │
└──────────────────────┬──────────────────────────┘
                       │ HTTP
                       ▼
┌──────────────────────────────────────────────────┐
│                  API SERVER (Go)                  │
│                                                   │
│  ┌─────────┐  ┌──────────┐  ┌────────────────┐  │
│  │ Router  │→ │ Handlers │→ │ Store (Postgres)│  │
│  └─────────┘  └──────────┘  └────────────────┘  │
│                       │                           │
│                       ▼                           │
│              ┌─────────────────┐                  │
│              │  Queue (Redis)  │                  │
│              └────────┬────────┘                  │
└───────────────────────┼──────────────────────────┘
                        │
           ┌────────────┼────────────┐
           ▼            ▼            ▼
      ┌────────┐  ┌────────┐  ┌────────┐
      │Worker 1│  │Worker 2│  │Worker N│
      └────────┘  └────────┘  └────────┘
```

## Components

### API Server

The entry point. A single Go binary that:

- Accepts HTTP requests
- Validates and stores data in Postgres
- Puts jobs into Redis when submitted
- Returns status and results

**Why Go?**
- Fast, compiled, single binary
- Great standard library (`net/http` is production-ready)
- Excellent concurrency with goroutines
- Strong ecosystem for infrastructure tools

**Why not Node.js/Python?**
They work, but Go gives us:
- Lower memory footprint (~10MB vs ~100MB+)
- No runtime dependency
- Better performance for CPU-bound work
- Compiled = no "it works on my machine" issues

### PostgreSQL

The source of truth. Stores:

- **Workflows**: What work needs to be done
- **Jobs**: Individual tasks within workflows
- **Job Events**: Audit trail of every state change
- **Workers**: Who's doing what

**Why Postgres?**
- ACID transactions (data integrity)
- JSONB support (flexible payloads)
- Excellent indexing for queries
- Battle-tested, runs everywhere

### Redis

The job queue. When a job is submitted:

1. API writes job to Postgres (source of truth)
2. API pushes job ID to Redis Stream (notification)
3. Worker reads from Redis Stream
4. Worker processes and updates Postgres

**Why Redis Streams?**
- Consumer groups = multiple workers, each job processed once
- Visibility timeout = if worker crashes, job reappears
- Atomic operations = no race conditions
- Fast (in-memory)

**Why not just Postgres polling?**
Polling works but wastes resources. Redis gives us:
- Real-time push to workers
- Built-in consumer groups
- Message persistence
- Better throughput under load

### Workers

Independent processes that:
1. Pull jobs from Redis
2. Execute the job logic
3. Update status in Postgres
4. Handle retries on failure

Workers can run:
- In the same process as the API (development)
- As separate processes (production)
- On different machines (distributed)

## Job Lifecycle

```
         ┌──────────┐
         │ CREATED  │  Job submitted via API
         └────┬─────┘
              │
              ▼
         ┌──────────┐
         │  QUEUED  │  Pushed to Redis Stream
         └────┬─────┘
              │
              ▼
         ┌──────────┐
         │PICKED_UP │  Worker claimed the job
         └────┬─────┘
              │
              ▼
         ┌──────────┐
         │ RUNNING  │  Worker executing
         └────┬─────┘
              │
       ┌──────┴──────┐
       ▼             ▼
┌──────────┐   ┌──────────┐
│COMPLETED │   │  FAILED  │
└──────────┘   └────┬─────┘
                    │
              ┌─────┴─────┐
              ▼           ▼
         ┌──────────┐ ┌──────────────┐
         │ RETRYING │ │ DEAD_LETTER  │
         └──────────┘ │ (max retries)│
                      └──────────────┘
```

Each state transition is recorded as a `job_event` in the database. This gives you a complete audit trail.

## Database Design

### Why UUIDs Instead of Auto-Increment?

```
Auto-increment: 1, 2, 3, 4, 5...
UUID:            550e8400-e29b-41d4-a716-446655440000
```

UUIDs are better for distributed systems because:
- No coordination needed between nodes
- Can generate IDs without database access
- No information leakage (can't guess next ID)
- Merge-safe (two DBs can merge without conflicts)

### Why JSONB for Payloads?

Instead of rigid columns:

```sql
-- Bad: rigid
CREATE TABLE jobs (
    image_url VARCHAR(255),
    width INTEGER,
    height INTEGER,
    format VARCHAR(10)
);

-- Good: flexible
CREATE TABLE jobs (
    payload JSONB  -- {"image": "photo.jpg", "width": 800}
);
```

Different job types need different data. JSONB lets us handle that without schema changes.

### Event Sourcing for Jobs

Every state change creates an event:

```sql
INSERT INTO job_events (job_id, event_type, data)
VALUES ('job-123', 'STATUS_CHANGED', '{"from": "QUEUED", "to": "RUNNING"}');
```

This gives you:
- Complete history of what happened
- Ability to replay/debug
- No data loss on status changes
- Audit trail for compliance

## Configuration

All configuration via environment variables (loaded from `.env` if present):

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | API server port |
| `READ_TIMEOUT_SECONDS` | 0 | Server read timeout |
| `WRITE_TIMEOUT_SECONDS` | 0 | Server write timeout |
| `IDLE_TIMEOUT_SECONDS` | 0 | Server idle timeout |
| `CORS_ALLOWED_ORIGINS` | *(none)* | Comma-separated list of allowed browser origins; CORS disabled when empty |
| `DB_HOST` | localhost | PostgreSQL host |
| `DB_PORT` | 5432 | PostgreSQL port |
| `DB_USER` | luthia | PostgreSQL user |
| `DB_PASSWORD` | luthia | PostgreSQL password |
| `DB_NAME` | luthia | PostgreSQL database |

:::note[Planned]

`REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB`, `WORKER_CONCURRENCY` and `WORKER_POLL_INTERVAL` are defined in the config struct but unused until the queue and worker land.

:::

## Project Structure

```
luthia/
├── cmd/                        # Entry points (each binary)
│   ├── api/main.go             # API server
│   ├── worker/main.go          # Worker process (planned)
│   └── migrate/main.go         # DB migration tool
├── internal/                   # Private code (can't import externally)
│   ├── api/                    # HTTP layer, router mount
│   │   └── handlers/           # Health handler
│   ├── config/                 # Configuration loading
│   ├── json/                   # JSON read/write helpers
│   ├── store/                  # Database access
│   │   └── postgres/           # pgx pool + migrations + sqlc
│   │       └── sqlc/           # Generated types & queries
│   └── workflows/              # Workflow service + handlers
├── bin/                        # Compiled binaries
├── docs/                       # Docusaurus documentation
├── docker-compose.yml          # Local development
├── Makefile
├── sqlc.yaml
└── PLAN.md                     # Development roadmap
```

**Why `internal/`?**
Go enforces that packages in `internal/` can't be imported by code outside the module. This keeps our implementation details private and prevents external coupling.

**How the layers talk:**
```
handlers (HTTP) → service (business rules) → sqlc (store.Querier) → pgx → Postgres
```

- `internal/workflows/service.go` defines a `Service` interface; handlers depend on the interface, not the implementation. This keeps HTTP code decoupled from storage and makes unit testing easy.
- `internal/store/postgres/sqlc` contains the generated types (`Workflow`, `CreateWorkflowParams`, …) and the `Querier` interface. Handlers and services import these types directly — there is no separate `model` package.
- Type mapping is configured in `sqlc.yaml`: UUIDs become `github.com/google/uuid.UUID`, `jsonb` becomes `encoding/json.RawMessage`, and timestamps become `time.Time`. DB enums generate typed constants with a `Valid()` method.

## Scaling Model

### V1: Single Process
```
API + Worker → Postgres + Redis
```

### V2: Separate Processes
```
API → Postgres + Redis ← Worker(s)
```

### V3: Distributed
```
API (multiple) → Postgres + Redis ← Workers (multiple machines)
                                       ↓
                                   Kubernetes
```

Each step adds complexity only when needed. Start simple, scale when necessary.

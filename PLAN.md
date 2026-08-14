# Luthia — Development Plan

> Open-source infrastructure for running reliable background workflows.

---

## Vision

A developer-first distributed workflow engine. Simple to start, powerful at scale.

```
YOUR APPLICATION
       │
       ▼
┌─────────────┐
│   LUTHIA    │
│             │
│  Workflows  │
│  Scheduler  │
│  Workers    │
│  Retries    │
│  State      │
└──────┬──────┘
       │
┌──────┼───────┐
▼      ▼       ▼
Worker Worker  Worker
```

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.22+ |
| Database | PostgreSQL 16 |
| Queue | Redis 7 (Streams) |
| Frontend | React 18 + TypeScript + Vite |
| Styling | Tailwind CSS |
| Container | Docker + Docker Compose |
| Docs | Docusaurus 3 |
| Future | Kubernetes, Helm |

---

## Architecture

### Core Concepts

- **Workflow**: A named sequence of steps (DAG or linear chain)
- **Job**: A single unit of work within a workflow
- **Worker**: A process that pulls jobs from a queue and executes them
- **Step**: A defined task type (e.g., "resize-image", "send-email")
- **Attempt**: Each execution try of a job (for retries)

### Job Lifecycle

```
CREATED → QUEUED → PICKED_UP → RUNNING → COMPLETED
                                ↓
                            FAILED → RETRYING → QUEUED (again)
                                ↓
                            DEAD_LETTER (max retries exceeded)
```

### Database Schema (core tables)

```
workflows     — id, name, description, config (JSONB), status, created_at
jobs          — id, workflow_id, step, payload (JSONB), status, priority,
                attempt_count, max_retries, scheduled_at, created_at, updated_at
job_events    — id, job_id, event_type, data (JSONB), created_at
workers       — id, hostname, status, current_job_id, last_heartbeat
```

### API (REST, v1)

```
POST   /api/workflows                    — Create workflow
GET    /api/workflows                    — List workflows
GET    /api/workflows/:id                — Get workflow details

POST   /api/workflows/:id/jobs           — Submit jobs to workflow
GET    /api/workflows/:id/jobs           — List jobs for workflow

GET    /api/jobs/:id                     — Get job detail
GET    /api/jobs/:id/events              — Get job event log
POST   /api/jobs/:id/retry               — Manually retry failed job
POST   /api/jobs/:id/cancel              — Cancel running job

GET    /api/workers                      — List active workers
GET    /api/stats                        — Dashboard stats
```

---

## Phased Development

### V1 — Foundation (Weeks 1-4)

**Goal**: A single user can create workflows, submit jobs, see them in a dashboard.
**Docs**: Introduction, Quickstart, basic API reference.

```
Week 1: Project scaffold
  - Go project structure (cmd/, internal/, pkg/)
  - Postgres schema + migrations (golang-migrate)
  - Docker Compose (Go app + Postgres + Redis)
  - Health check endpoint

Week 2: Core API
  - CRUD workflows
  - Submit jobs (store in DB with status QUEUED)
  - List/get jobs
  - Job status transitions (state machine)

Week 3: Basic worker
  - Go worker process that polls Postgres for QUEUED jobs
  - Picks up job → marks RUNNING → executes → marks COMPLETED
  - Simulated step handlers (sleep, process, fail randomly)
  - Basic retry on failure (max 3 attempts)

Week 4: Dashboard v1
  - React + Vite + TypeScript scaffold
  - Dashboard page: running/completed/failed counts
  - Workflow list page
  - Job list page per workflow
  - Job detail page with attempt history
  - Docker Compose includes frontend

Week 4 (docs): Docusaurus scaffold
  - Docusaurus 3 project in /docs
  - "What is Luthia?" intro page
  - Quickstart guide (docker compose up → first workflow)
  - Basic API reference (endpoints list)
  - docs:dev script running alongside the app
```

**Deliverable**: `docker compose up` → full working system on localhost.

---

### V2 — Core Engine (Weeks 5-10)

**Goal**: Production-grade job processing with queues, retries, and idempotency.
**Docs**: Workers guide, retry configuration, queue concepts.

```
Week 5-6: Redis queue
  - Replace Postgres polling with Redis Streams
  - Producer: API enqueues jobs
  - Consumer: Worker reads from stream
  - Consumer groups for multiple workers
  - Visibility timeout (if worker crashes, job reappears)

Week 7-8: Retry & idempotency
  - Exponential backoff retry strategy
  - Configurable max retries per job
  - Dead letter queue for permanently failed jobs
  - Idempotency keys (prevent duplicate execution)
  - Job deduplication

Week 9-10: Worker management
  - Worker registration (heartbeat to DB)
  - Worker status tracking (idle/busy/offline)
  - Graceful shutdown (finish current job, stop accepting)
  - Concurrent job processing per worker (configurable)
  - Job timeout (kill stuck jobs)
```

---

### V3 — Observability & Scheduling (Weeks 11-18)

**Goal**: See everything, schedule anything.
**Docs**: Scheduling, metrics, observability guide.

```
Week 11-12: Job events & audit log
  - Event sourcing for job lifecycle
  - Every state change emits an event
  - Events stored in job_events table
  - API to query job history/timeline

Week 13-14: Scheduling
  - Delayed jobs (execute after N seconds)
  - Cron-based recurring workflows
  - Schedule management API
  - Next-run calculation

Week 15-16: Metrics & tracing
  - Prometheus metrics endpoint
  - Job duration, throughput, error rates
  - Per-step performance stats
  - Request tracing with correlation IDs
  - Structured logging (slog)

Week 17-18: Dashboard v2
  - Real-time updates (WebSocket/SSE)
  - Live job progress
  - Worker status panel
  - Metrics charts (Grafana integration)
  - Job event timeline view
  - Search and filter jobs
```

---

### V4 — Distribution & Scale (Weeks 19-26)

**Goal**: Run across multiple machines.
**Docs**: Kubernetes deployment, distributed setup, multi-tenancy.

```
Week 19-20: Distributed workers
  - Workers run on separate machines
  - Service discovery (static config → Consul/etcd)
  - Network partitioning handling
  - Split-brain prevention

Week 21-22: Horizontal scaling
  - Multiple API instances behind load balancer
  - Sticky sessions not needed (stateless API)
  - Database connection pooling (pgxpool)
  - Redis connection management

Week 23-24: Kubernetes
  - Helm chart for Luthia
  - Worker Deployment with HPA
  - ConfigMaps for configuration
  - Persistent volume for Postgres
  - Health probes (liveness/readiness)

Week 25-26: Multi-tenancy
  - Workspace/namespace isolation
  - API keys and authentication
  - Rate limiting
  - Usage tracking per tenant
```

---

### V5 — Intelligence & Ecosystem (Weeks 27-34)

**Goal**: AI-assisted operations and community.
**Docs**: SDK reference, MCP integration, contributing guide, full docs site polish.

```
Week 27-28: MCP Server
  - Tool: query_jobs(filter) → list jobs
  - Tool: get_job_status(id) → job detail
  - Tool: retry_job(id) → trigger retry
  - Tool: workflow_stats(id) → performance metrics
  - Resource: recent_failures → last 24h failures

Week 29-30: SDK
  - Go SDK: client.Submit("workflow-name", payload)
  - TypeScript SDK: client.workflow("name").submit(data)
  - Retry configuration
  - Error handling

Week 31-32: Workflow builder
  - Define workflows via YAML/JSON config
  - DAG support (parallel branches, joins)
  - Conditional steps
  - Step dependencies

Week 33-34: Polish & open source
  - Documentation site (Docusaurus)
  - Contributing guide
  - Issue templates
  - CI/CD (GitHub Actions)
  - License (MIT or Apache 2.0)
  - README with quickstart
  - Logo and branding
```

---

## Project Structure

```
luthia/
├── cmd/
│   ├── api/              # API server entrypoint
│   ├── worker/           # Worker entrypoint
│   └── migrate/          # Database migration tool
├── internal/
│   ├── api/              # HTTP handlers, router
│   ├── api/handlers/     # Route handlers
│   ├── api/middleware/    # Auth, logging, etc.
│   ├── engine/           # Core workflow engine
│   ├── engine/state/     # Job state machine
│   ├── queue/            # Queue abstraction (Redis impl)
│   ├── store/            # Database access layer
│   ├── store/postgres/   # Postgres implementation
│   ├── worker/           # Worker logic
│   ├── scheduler/        # Cron/delayed job scheduling
│   └── model/            # Domain types
├── pkg/
│   └── sdk/              # Public Go SDK
├── dashboard/            # React frontend
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── hooks/
│   │   └── api/          # API client
│   └── package.json
├── migrations/           # SQL migrations
├── deploy/               # Docker, Kubernetes configs
├── docs/                 # Docusaurus documentation site
│   ├── docs/             # Markdown content
│   ├── src/              # Docusaurus components
│   ├── docusaurus.config.js
│   └── package.json
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
└── go.sum
```

---

## Documentation (Docusaurus)

Documentation is built incrementally alongside development. Not at the end. Each version ships with its docs.

### Structure

```
docs/
├── introduction.md
├── quickstart.md
├── concepts/
│   ├── workflows.md
│   ├── jobs.md
│   ├── workers.md
│   └── events.md
├── guides/
│   ├── creating-workflows.md
│   ├── submitting-jobs.md
│   ├── configuring-retries.md
│   └── deploying-with-docker.md
├── api/
│   ├── rest-reference.md
│   └── sdks.md
├── reference/
│   ├── configuration.md
│   ├── cli.md
│   └── architecture.md
└── contributing.md
```

### Per-Version Documentation

| Version | Docs delivered |
|---------|---------------|
| V1 | Introduction, Quickstart, "What is Luthia", basic API reference |
| V2 | Workers guide, retry configuration, queue concepts |
| V3 | Scheduling, metrics, observability guide |
| V4 | Kubernetes deployment, distributed setup, multi-tenancy |
| V5 | SDK reference, MCP integration, contributing guide |

### Docusaurus Setup

- Lives in `/docs` at project root
- `package.json` scripts: `docs:dev`, `docs:build`, `docs:deploy`
- Deployed to GitHub Pages or Vercel on release
- Versioned docs (Docusaurus built-in versioning)
- Search with local search plugin (no Algolia needed initially)

---

## Naming & Positioning

**Luthia** — from "relay" concept. Workflows flow through the system like signals through a relay.

**Tagline options**:
- "Run, monitor and scale background workflows."
- "Open-source infrastructure for reliable background jobs."
- "The workflow engine that gets out of your way."

**Comparable projects** (for positioning, not competition):
| Project | Focus |
|---------|-------|
| Temporal | Complex workflows, heavy |
| BullMQ | Node.js ecosystem |
| Celery | Python ecosystem |
| Inngest | Serverless-first |
| **Luthia** | Simple, Go-native, self-hosted |

---

## How to Present to Companies

**Totalis** (fintech):
> "Distributed workers, concurrency control, idempotency, queue management, failure recovery — all core banking infrastructure concepts."

**Dedalus** (healthcare):
> "Scheduling, fault tolerance, orchestration, Kubernetes-native — exactly what clinical systems need."

**Nessie** (AI agents):
> "Go backend, clean APIs, extensible architecture — perfect for agent task orchestration."

---

## Immediate Next Steps

1. Initialize Go module
2. Set up Docker Compose (Go + Postgres + Redis)
3. Create database migrations
4. Implement health check endpoint
5. Build basic API (CRUD workflows + submit jobs)
6. Build worker that polls and processes jobs
7. Scaffold React dashboard

**Start with**: `go mod init github.com/DecoXFE/luthia`

**Day 1 includes**: Docusaurus scaffold alongside Go project — docs start from minute zero.

# Luthia

> Open-source infrastructure for reliable background workflows.

## Work in Progress

Luthia is under active development and is not ready for production use. APIs, configuration, and architecture may change while the core workflow engine is being built.

## About

Luthia is a Go-based, self-hosted workflow engine for defining workflows, submitting background jobs, and tracking their execution. It is designed to provide reliable job processing with retries, state tracking, and eventually distributed workers, scheduling, and observability.

## Current Status

Implemented so far:

- Workflow CRUD API
- PostgreSQL storage, migrations, and typed queries with sqlc
- Job state transitions and runtime tracking
- Health check endpoint
- Docker Compose environment with PostgreSQL, Redis, API, documentation, and Redis Commander

Planned next: queue-backed workers, retries, scheduling, dashboard, and distributed execution.

## Stack

- Go
- PostgreSQL
- Redis
- Docker Compose
- Docusaurus documentation

## Getting Started

Start the local environment with Docker:

```bash
make dev
```

The API is available at `http://localhost:8080` and the documentation at `http://localhost:3000`.

To run the test suite:

```bash
make test
```

For more information you can use:
```bash
make help
```

See [PLAN.md](PLAN.md) for the development roadmap.



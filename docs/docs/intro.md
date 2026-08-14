# Introduction

**Luthia** is an open-source workflow engine for running reliable background jobs.

If you've ever needed to process thousands of images, scrape thousands of URLs, run ML experiments, or batch-process documents — you know the problem:

> "I have a lot of work to do, and I need to do it reliably."

Luthia solves that.

## What Luthia Does

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

You define **workflows** (sequences of steps), submit **jobs** (units of work), and Luthia handles:

- **Distribution** across multiple workers
- **Retries** when things fail
- **State tracking** so you know what's happening
- **Observability** via a dashboard and API

## Who It's For

Developers who need to:

- Process large batches of data
- Run background tasks reliably
- Scale job processing horizontally
- Monitor what's happening with their jobs

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go |
| Database | PostgreSQL |
| Queue | Redis ( Streams) |
| Frontend | React + TypeScript |
| Docs | Docusaurus |

## Quick Links

- [Quickstart](quickstart) — Get running in 5 minutes
- [Architecture](architecture) — How Luthia works internally
- [Concepts](concepts/workflows) — Core building blocks
- [API Reference](api/rest-reference) — REST endpoint documentation

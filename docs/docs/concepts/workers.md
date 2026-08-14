# Workers

A **worker** is a process that pulls jobs from the queue and executes them.

## How Workers Work

```
1. Worker starts, registers itself
2. Polls Redis Stream for available jobs
3. Claims a job (marks as PICKED_UP)
4. Executes the job logic
5. Updates status (COMPLETED or FAILED)
6. Repeats
```

## Worker Lifecycle

```
         ┌──────────┐
         │  IDLE    │  No job assigned
         └────┬─────┘
              │
              ▼
         ┌──────────┐
         │  BUSY    │  Executing a job
         └────┬─────┘
              │
              ▼
         ┌──────────┐
         │   IDLE   │  Ready for next job
         └──────────┘
```

## Concurrency

Each worker can process multiple jobs concurrently:

```
Worker (concurrency: 5)
├── Job 1 (running)
├── Job 2 (running)
├── Job 3 (running)
├── Job 4 (running)
└── Job 5 (running)
```

Configure with `WORKER_CONCURRENCY` environment variable.

## Failure Handling

### Worker Crashes Mid-Job

If a worker dies while processing:

1. Job stays in "RUNNING" state
2. Heartbeat stops
3. After timeout, job re-enters queue
4. Another worker picks it up

### Job Timeout

If a job runs too long:

```
Max execution time: 300 seconds
Elapsed: 301 seconds
→ Job killed, marked as FAILED
→ Retried on another worker
```

## Graceful Shutdown

When you stop a worker:

1. Stops accepting new jobs
2. Finishes currently running jobs
3. Updates status to idle
4. Process exits cleanly

This prevents:
- Lost jobs
- Corrupted state
- Partial execution

## Worker Registration

Each worker registers in the database:

```json
{
  "id": "worker-abc",
  "hostname": "server-1",
  "status": "idle",
  "last_heartbeat": "2025-01-15T10:40:00Z"
}
```

The API can then show:
- How many workers are active
- Which worker is processing which job
- Worker health and uptime

## Multiple Workers

```
┌─────────────────────────────────────────┐
│              Redis Stream               │
│  [job1] [job2] [job3] [job4] [job5]    │
└─────────────┬──────────┬────────────────┘
              │          │
     ┌────────┴──┐  ┌───┴────────┐
     │ Worker 1  │  │  Worker 2  │
     │ (idle)    │  │  (busy)    │
     └───────────┘  └────────────┘
```

- Each job processed by exactly one worker (consumer group)
- If worker dies, job returns to queue
- Workers can run on same or different machines

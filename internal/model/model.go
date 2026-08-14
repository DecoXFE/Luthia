package model

import "time"

type JobStatus string

const (
	JobStatusCreated   JobStatus = "CREATED"
	JobStatusQueued    JobStatus = "QUEUED"
	JobStatusPickedUp  JobStatus = "PICKED_UP"
	JobStatusRunning   JobStatus = "RUNNING"
	JobStatusCompleted JobStatus = "COMPLETED"
	JobStatusFailed    JobStatus = "FAILED"
	JobStatusRetrying  JobStatus = "RETRYING"
	JobStatusDead      JobStatus = "DEAD_LETTER"
	JobStatusCancelled JobStatus = "CANCELLED"
)

type Workflow struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Job struct {
	ID           string    `json:"id"`
	WorkflowID   string    `json:"workflow_id"`
	Step         string    `json:"step"`
	Payload      []byte    `json:"payload"`
	Status       JobStatus `json:"status"`
	Priority     int       `json:"priority"`
	AttemptCount int       `json:"attempt_count"`
	MaxRetries   int       `json:"max_retries"`
	ScheduledAt  *time.Time `json:"scheduled_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type JobEvent struct {
	ID        string    `json:"id"`
	JobID     string    `json:"job_id"`
	EventType string    `json:"event_type"`
	Data      []byte    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

type Worker struct {
	ID           string     `json:"id"`
	Hostname     string     `json:"hostname"`
	Status       string     `json:"status"`
	CurrentJobID *string    `json:"current_job_id,omitempty"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
}

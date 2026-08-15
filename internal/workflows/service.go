package workflows

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	store "github.com/DecoXFE/luthia/internal/store/postgres/sqlc"
)

const (
	maxNameLength  = 255
	pgErrCodeUniqueViolation = "23505"
)

var (
	ErrNotFound    = errors.New("workflow not found")
	ErrNameTaken   = errors.New("workflow name already exists")
	ErrInvalidName = errors.New("workflow name is required")
)

type Service interface {
	List(ctx context.Context) ([]store.Workflow, error)
	Create(ctx context.Context, tempWorkflow store.CreateWorkflowParams) (store.Workflow, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	queries store.Querier
}

func NewService(queries store.Querier) Service {
	return &service{queries: queries}
}

func (s *service) List(ctx context.Context) ([]store.Workflow, error) {
	return s.queries.ListWorkflows(ctx)
}

func (s *service) Create(ctx context.Context, tempWorkflow store.CreateWorkflowParams) (store.Workflow, error) {
	if strings.TrimSpace(tempWorkflow.Name) == "" {
		return store.Workflow{}, ErrInvalidName
	}
	if len(tempWorkflow.Name) > maxNameLength {
		return store.Workflow{}, fmt.Errorf("%w: name exceeds %d characters", ErrInvalidName, maxNameLength)
	}

	workflow, err := s.queries.CreateWorkflow(ctx, tempWorkflow)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgErrCodeUniqueViolation {
			return store.Workflow{}, ErrNameTaken
		}

		return store.Workflow{}, err
	}

	return workflow, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.queries.DeleteWorkflow(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	return err
}

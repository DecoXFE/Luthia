package workflows

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	store "github.com/DecoXFE/luthia/internal/store/postgres/sqlc"
)

type MockQuerier struct {
	mock.Mock
	store.Querier
}

func (m *MockQuerier) CreateWorkflow(ctx context.Context, arg store.CreateWorkflowParams) (store.Workflow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(store.Workflow), args.Error(1)
}

func (m *MockQuerier) DeleteWorkflow(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockQuerier) ListWorkflows(ctx context.Context) ([]store.Workflow, error) {
	args := m.Called(ctx)
	return args.Get(0).([]store.Workflow), args.Error(1)
}

func TestServiceCreateSuccess(t *testing.T) {
	// Arrange
	want := store.Workflow{ID: uuid.New(), Name: "test", Status: store.WorkflowStatusActive}
	querier := &MockQuerier{}
	querier.On("CreateWorkflow", mock.Anything, mock.Anything).Return(want, nil)
	service := NewService(querier)

	// Act
	got, err := service.Create(context.Background(), store.CreateWorkflowParams{Name: "test", Description: new("desc")})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, want, got)
	querier.AssertExpectations(t)

	params := querier.Calls[0].Arguments.Get(1).(store.CreateWorkflowParams)
	assert.Equal(t, "test", params.Name)
	assert.Equal(t, new("desc"), params.Description)
}

func TestServiceCreateInvalidName(t *testing.T) {
	// Arrange
	tests := []struct {
		name string
		in   string
	}{
		{name: "empty", in: ""},
		{name: "whitespace only", in: "   "},
		{name: "too long", in: strings.Repeat("a", maxNameLength+1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &MockQuerier{}
			s := NewService(q)

			// Act
			_, err := s.Create(context.Background(), store.CreateWorkflowParams{Name: tt.in})

			// Assert
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrInvalidName)
			q.AssertNotCalled(t, "CreateWorkflow", mock.Anything, mock.Anything)
		})
	}
}

func TestServiceCreateNameTaken(t *testing.T) {
	// Arrange
	q := &MockQuerier{}
	q.On("CreateWorkflow", mock.Anything, mock.Anything).Return(store.Workflow{}, &pgconn.PgError{Code: pgErrCodeUniqueViolation})
	s := NewService(q)

	// Act
	_, err := s.Create(context.Background(), store.CreateWorkflowParams{Name: "dup"})

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNameTaken)
	q.AssertExpectations(t)
}

func TestServiceCreateQueryError(t *testing.T) {
	// Arrange
	boom := errors.New("db down")
	q := &MockQuerier{}
	q.On("CreateWorkflow", mock.Anything, mock.Anything).Return(store.Workflow{}, boom)
	s := NewService(q)

	// Act
	_, err := s.Create(context.Background(), store.CreateWorkflowParams{Name: "test"})

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, boom)
	q.AssertExpectations(t)
}

func TestServiceList(t *testing.T) {
	// Arrange
	want := []store.Workflow{{ID: uuid.New(), Name: "a"}, {ID: uuid.New(), Name: "b"}}
	q := &MockQuerier{}
	q.On("ListWorkflows", mock.Anything).Return(want, nil)
	s := NewService(q)

	// Act
	got, err := s.List(context.Background())

	// Assert
	require.NoError(t, err)
	assert.Equal(t, want, got)
	q.AssertExpectations(t)
}

func TestServiceDeleteSuccess(t *testing.T) {
	// Arrange
	id := uuid.New()
	q := &MockQuerier{}
	q.On("DeleteWorkflow", mock.Anything, id).Return(id, nil)
	s := NewService(q)

	// Act
	err := s.Delete(context.Background(), id)

	// Assert
	require.NoError(t, err)
	q.AssertExpectations(t)
}

func TestServiceDeleteNotFound(t *testing.T) {
	// Arrange
	q := &MockQuerier{}
	q.On("DeleteWorkflow", mock.Anything, mock.Anything).Return(uuid.Nil, pgx.ErrNoRows)
	s := NewService(q)

	// Act
	err := s.Delete(context.Background(), uuid.New())

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
	q.AssertExpectations(t)
}

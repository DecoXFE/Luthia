package workflows

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	store "github.com/DecoXFE/luthia/internal/store/postgres/sqlc"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Create(ctx context.Context, params store.CreateWorkflowParams) (store.Workflow, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(store.Workflow), args.Error(1)
}

func (m *MockService) List(ctx context.Context) ([]store.Workflow, error) {
	args := m.Called(ctx)
	return args.Get(0).([]store.Workflow), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func newTestHandler(s Service) http.Handler {
	r := chi.NewRouter()
	NewHandler(s).Enroute(r)
	return r
}

func TestHandlerCreate(t *testing.T) {
	// Arrange
	want := store.Workflow{ID: uuid.New(), Name: "test", Status: store.WorkflowStatusActive}
	s := &MockService{}
	s.On("Create", mock.Anything, mock.Anything).Return(want, nil)
	h := newTestHandler(s)

	req := httptest.NewRequest("POST", "/api/workflows", strings.NewReader(`{"name":"test","description":"desc"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	h.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusCreated, w.Code)
	s.AssertExpectations(t)

	params := s.Calls[0].Arguments.Get(1).(store.CreateWorkflowParams)
	assert.Equal(t, "test", params.Name)
	assert.Equal(t, new("desc"), params.Description)
	assert.JSONEq(t, `{"id":"`+want.ID.String()+`","name":"test","description":null,"status":"active","config":null,"created_at":null,"updated_at":null}`, w.Body.String())
}

func TestHandlerCreateInvalidBody(t *testing.T) {
	// Arrange
	s := &MockService{}
	h := newTestHandler(s)

	req := httptest.NewRequest("POST", "/api/workflows", strings.NewReader(`{not json`))
	w := httptest.NewRecorder()

	// Act
	h.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"invalid request body"}`, w.Body.String())
	s.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestHandlerCreateUnknownField(t *testing.T) {
	// Arrange
	s := &MockService{}
	h := newTestHandler(s)

	req := httptest.NewRequest("POST", "/api/workflows", strings.NewReader(`{"name":"test","bogus":1}`))
	w := httptest.NewRecorder()

	// Act
	h.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"invalid request body"}`, w.Body.String())
	s.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestHandlerCreateInvalidName(t *testing.T) {
	// Arrange
	s := &MockService{}
	s.On("Create", mock.Anything, mock.Anything).Return(store.Workflow{}, ErrInvalidName)
	h := newTestHandler(s)

	req := httptest.NewRequest("POST", "/api/workflows", strings.NewReader(`{"name":""}`))
	w := httptest.NewRecorder()

	// Act
	h.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"workflow name is required"}`, w.Body.String())
	s.AssertExpectations(t)
}

func TestHandlerCreateNameTaken(t *testing.T) {
	// Arrange
	s := &MockService{}
	s.On("Create", mock.Anything, mock.Anything).Return(store.Workflow{}, ErrNameTaken)
	h := newTestHandler(s)

	req := httptest.NewRequest("POST", "/api/workflows", strings.NewReader(`{"name":"dup"}`))
	w := httptest.NewRecorder()

	// Act
	h.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusConflict, w.Code)
	assert.JSONEq(t, `{"error":"workflow name already exists"}`, w.Body.String())
	s.AssertExpectations(t)
}

func TestHandlerCreateServiceError(t *testing.T) {
	// Arrange
	s := &MockService{}
	s.On("Create", mock.Anything, mock.Anything).Return(store.Workflow{}, errors.New("db down"))
	h := newTestHandler(s)

	req := httptest.NewRequest("POST", "/api/workflows", strings.NewReader(`{"name":"test"}`))
	w := httptest.NewRecorder()

	// Act
	h.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusInternalServerError, w.Code)
	s.AssertExpectations(t)
}

func TestHandlerList(t *testing.T) {
	// Arrange
	want := []store.Workflow{{ID: uuid.New(), Name: "a"}, {ID: uuid.New(), Name: "b"}}
	s := &MockService{}
	s.On("List", mock.Anything).Return(want, nil)
	h := newTestHandler(s)

	req := httptest.NewRequest("GET", "/api/workflows", nil)
	w := httptest.NewRecorder()

	// Act
	h.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `[{"id":"`+want[0].ID.String()+`","name":"a","description":null,"status":"","config":null,"created_at":null,"updated_at":null},{"id":"`+want[1].ID.String()+`","name":"b","description":null,"status":"","config":null,"created_at":null,"updated_at":null}]`, w.Body.String())
	s.AssertExpectations(t)
}

func TestHandlerListServiceError(t *testing.T) {
	// Arrange
	s := &MockService{}
	s.On("List", mock.Anything).Return([]store.Workflow(nil), errors.New("db down"))
	h := newTestHandler(s)

	req := httptest.NewRequest("GET", "/api/workflows", nil)
	w := httptest.NewRecorder()

	// Act
	h.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusInternalServerError, w.Code)
	s.AssertExpectations(t)
}

func TestHandlerDelete(t *testing.T) {
	// Arrange
	id := uuid.New()
	s := &MockService{}
	s.On("Delete", mock.Anything, id).Return(nil)
	h := newTestHandler(s)

	req := httptest.NewRequest("DELETE", "/api/workflows/"+id.String(), nil)
	w := httptest.NewRecorder()

	// Act
	h.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusNoContent, w.Code)
	s.AssertExpectations(t)
}

func TestHandlerDeleteInvalidID(t *testing.T) {
	// Arrange
	s := &MockService{}
	h := newTestHandler(s)

	req := httptest.NewRequest("DELETE", "/api/workflows/not-a-uuid", nil)
	w := httptest.NewRecorder()

	// Act
	h.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"invalid workflow id"}`, w.Body.String())
	s.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

func TestHandlerDeleteNotFound(t *testing.T) {
	// Arrange
	s := &MockService{}
	s.On("Delete", mock.Anything, mock.Anything).Return(ErrNotFound)
	h := newTestHandler(s)

	req := httptest.NewRequest("DELETE", "/api/workflows/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()

	// Act
	h.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"error":"workflow not found"}`, w.Body.String())
	s.AssertExpectations(t)
}

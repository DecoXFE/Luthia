package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer_HealthRoute(t *testing.T) {
	// Arrange
	mux := NewServer()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Act
	mux.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}

func TestNewServer_UnknownRoute(t *testing.T) {
	// Arrange
	mux := NewServer()

	req := httptest.NewRequest("GET", "/does-not-exist", nil)
	w := httptest.NewRecorder()

	// Act
	mux.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusNotFound, w.Code)
}

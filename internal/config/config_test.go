package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_WithEnvVars(t *testing.T) {
	// Arrange
	os.Setenv("PORT", "9090")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5433")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")

	// Act
	cfg := Load()

	// Assert
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "db.example.com", cfg.Database.Host)
	assert.Equal(t, 5433, cfg.Database.Port)
}

func TestLoad_EmptyEnv(t *testing.T) {
	// Arrange
	os.Unsetenv("PORT")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("REDIS_ADDR")
	os.Unsetenv("REDIS_PASSWORD")
	os.Unsetenv("REDIS_DB")
	os.Unsetenv("WORKER_CONCURRENCY")
	os.Unsetenv("WORKER_POLL_INTERVAL")

	// Act
	cfg := Load()

	// Assert
	assert.Equal(t, 0, cfg.Server.Port)
	assert.Empty(t, cfg.Database.Host)
}

func TestDSN(t *testing.T) {
	// Arrange
	db := DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "luthia",
		Password: "secret",
		Name:     "luthia",
	}

	// Act
	dsn := db.DSN()

	// Assert
	require.Equal(t, "postgres://luthia:secret@localhost:5432/luthia?sslmode=disable", dsn)
}

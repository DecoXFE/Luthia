package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/DecoXFE/luthia/internal/api"
	"github.com/DecoXFE/luthia/internal/config"
	"github.com/DecoXFE/luthia/internal/store/postgres"
)

// @title Luthia API
// @version 1.0
// @description Open-source infrastructure for running reliable background workflows.
// @BasePath /

// @tag.name Workflows
// @tag.description Workflows are the top-level units of work in Luthia. They group related jobs, define how they should run (the workflow "config") and track their overall state. The workflows API is the entry point: create a workflow to start tracking work, list them to see what exists, and delete them when they are no longer needed (removing their jobs too). To learn what a workflow is and the states it can be in, see the workflows concept guide: https://luthia.dev/docs/concepts/workflows
// @tag.docs.url https://luthia.dev/docs/concepts/workflows
// @tag.docs.description Understand what a workflow is, its lifecycle and its states.

// @tag.name System
// @tag.description Operational endpoints for checking that the API is alive and healthy. These are meant for load balancers, uptime probes and CI pipelines, not for application code. If you are integrating Luthia into your infrastructure, wire the health endpoint into your monitoring rather than poking at business endpoints.
func main() {
	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Configuration
	cfg := config.Load()

	// Database
	ctx := context.Background()

	db, err := postgres.New(ctx, cfg.Database.DSN())
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to database successfully")

	defer db.Close()

	// Startup
	api := api.Application{
		Config: *cfg,
		DbPool: db.Pool,
	}

	if err := api.Run(api.Mount()); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}

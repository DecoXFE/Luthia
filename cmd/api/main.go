package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/DecoXFE/luthia/internal/api"
	"github.com/DecoXFE/luthia/internal/config"
	"github.com/DecoXFE/luthia/internal/store/postgres"
)

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

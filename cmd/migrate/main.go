package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/DecoXFE/luthia/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: migrate <command>")
		fmt.Fprintln(os.Stderr, "commands: up, down, version")
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg := config.Load()

	m, err := migrate.New("file://internal/store/postgres/migrations", cfg.Database.DSN())
	if err != nil {
		logger.Error("failed to create migrator", "error", err)
		os.Exit(1)
	}
	defer m.Close()

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			logger.Error("migration failed", "error", err)
			os.Exit(1)
		}

		logger.Info("migrations applied")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			logger.Error("migration rollback failed", "error", err)
			os.Exit(1)
		}

		logger.Info("migrations rolled back")

	case "force":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: migrate force <version>")
			os.Exit(1)
		}

		version, err := strconv.Atoi(os.Args[2])

		if err != nil {
			logger.Error("invalid version", "error", err)
			os.Exit(1)
		}
		if err := m.Force(version); err != nil {
			logger.Error("failed to force version", "error", err)
			os.Exit(1)
		}

		logger.Info("forced version", "version", version)

	case "version":
		version, dirty, err := m.Version()

		if err != nil {
			logger.Error("failed to get version", "error", err)
			os.Exit(1)
		}

		logger.Info("current version", "version", version, "dirty", dirty)

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		fmt.Fprintln(os.Stderr, "commands: up, down, version")
		os.Exit(1)
	}
}

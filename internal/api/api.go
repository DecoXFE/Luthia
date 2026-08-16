package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/DecoXFE/luthia/internal/api/handlers"
	"github.com/DecoXFE/luthia/internal/config"
	store "github.com/DecoXFE/luthia/internal/store/postgres/sqlc"
	"github.com/DecoXFE/luthia/internal/workflows"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	config.Config
	DbPool *pgxpool.Pool
}

func (app *Application) Mount() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.ClientIPFromRemoteAddr)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	if len(app.Server.AllowedOrigins) > 0 {
		router.Use(cors.Handler(cors.Options{
			AllowedOrigins:   app.Server.AllowedOrigins,
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
			AllowCredentials: false,
			MaxAge:           300,
		}))
	}

	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/health", handlers.Health)

	workflowsService := workflows.NewService(store.New(app.DbPool))
	workflows.NewHandler(workflowsService).Enroute(router)

	return router
}

func (app *Application) Run(h http.Handler) error {
	addr := fmt.Sprintf(":%d", app.Server.Port)
	slog.Info("starting luthia api", "addr", addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      h,
		WriteTimeout: time.Second * time.Duration(app.Server.WriteTimeoutSecs),
		ReadTimeout:  time.Second * time.Duration(app.Server.ReadTimeoutSecs),
		IdleTimeout:  time.Second * time.Duration(app.Server.IdleTimeoutSecs),
	}

	return server.ListenAndServe()
}

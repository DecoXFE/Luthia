package api

import (
	"net/http"

	"github.com/DecoXFE/luthia/internal/api/handlers"
)

func NewServer() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.Health)

	return mux
}

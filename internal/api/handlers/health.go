package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

// Health godoc
// @Summary Health check
// @Description Returns the service health status.
// @Tags System
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func Health(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(HealthResponse{Status: "ok"})
}

package workflows

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/DecoXFE/luthia/internal/json"
	store "github.com/DecoXFE/luthia/internal/store/postgres/sqlc"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Enroute(r chi.Router) {
	r.Post("/api/workflows", h.Create)
	r.Get("/api/workflows", h.List)
	r.Delete("/api/workflows/{id}", h.Delete)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var tempWorkflow store.CreateWorkflowParams
	if err := json.Read(r, &tempWorkflow); err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	workflow, err := h.service.Create(r.Context(), tempWorkflow)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidName):
			json.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrNameTaken):
			json.WriteError(w, http.StatusConflict, err.Error())
		default:
			json.WriteError(w, http.StatusInternalServerError, "failed to create workflow")
		}
		return
	}

	json.Write(w, http.StatusCreated, workflow)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	workflows, err := h.service.List(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, "failed to list workflows")
		return
	}

	json.Write(w, http.StatusOK, workflows)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid workflow id")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			json.WriteError(w, http.StatusNotFound, err.Error())
			return
		}

		json.WriteError(w, http.StatusInternalServerError, "failed to delete workflow")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

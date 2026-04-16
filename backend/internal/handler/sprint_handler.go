package handler

import (
	"net/http"

	"github.com/floqast/task-management/backend/internal/middleware"
)

type SprintHandler struct{}

func NewSprintHandler() *SprintHandler { return &SprintHandler{} }

func (h *SprintHandler) ListSprints(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *SprintHandler) CreateSprint(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *SprintHandler) GetSprint(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *SprintHandler) UpdateSprint(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *SprintHandler) DeleteSprint(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *SprintHandler) AddTaskToSprint(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *SprintHandler) RemoveTaskFromSprint(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

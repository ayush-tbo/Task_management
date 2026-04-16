package handler

import (
	"net/http"

	"github.com/floqast/task-management/backend/internal/middleware"
)

type ActivityHandler struct{}

func NewActivityHandler() *ActivityHandler { return &ActivityHandler{} }

func (h *ActivityHandler) GetProjectActivity(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *ActivityHandler) GetTaskActivity(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

package handler

import (
	"net/http"

	"github.com/floqast/task-management/backend/internal/middleware"
)

type LabelHandler struct{}

func NewLabelHandler() *LabelHandler { return &LabelHandler{} }

func (h *LabelHandler) ListLabels(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *LabelHandler) CreateLabel(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *LabelHandler) UpdateLabel(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *LabelHandler) DeleteLabel(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

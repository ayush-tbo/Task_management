package handler

import "net/http"

type ActivityHandler struct{}

func NewActivityHandler() *ActivityHandler { return &ActivityHandler{} }

func (h *ActivityHandler) GetProjectActivity(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *ActivityHandler) GetTaskActivity(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

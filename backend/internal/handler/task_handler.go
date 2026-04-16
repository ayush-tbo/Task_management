package handler

import (
	"net/http"

	"github.com/floqast/task-management/backend/internal/middleware"
)

type TaskHandler struct{}

func NewTaskHandler() *TaskHandler { return &TaskHandler{} }

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *TaskHandler) AssignTask(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *TaskHandler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *TaskHandler) GetTaskTimeTracking(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *TaskHandler) LogTaskTime(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *TaskHandler) GetMyTasks(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

package handler

import (
	"log/slog"
	"net/http"

	"github.com/floqast/task-management/backend/internal/middleware"
	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/service"
)

type ActivityHandler struct {
	activityService *service.ActivityService
	logger          *slog.Logger
}

func NewActivityHandler(activityService *service.ActivityService, logger *slog.Logger) *ActivityHandler {
	return &ActivityHandler{
		activityService: activityService,
		logger:          logger,
	}
}

func (h *ActivityHandler) GetProjectActivity(w http.ResponseWriter, r *http.Request) {
	projectID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid project id param", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid project id")
		return
	}

	activities, err := h.activityService.FindByProject(r.Context(), projectID)
	if err != nil {
		h.logger.Error("get project activity failed", "error", err, "project_id", projectID)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve activities")
		return
	}

	if len(activities) == 0 {
		middleware.WriteJSON(w, http.StatusOK, map[string]any{"activities": []model.ActivityEntry{}})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, map[string]any{"activities": activities})
}

func (h *ActivityHandler) GetTaskActivity(w http.ResponseWriter, r *http.Request) {
	taskID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid task id param", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	activities, err := h.activityService.FindByTask(r.Context(), taskID)
	if err != nil {
		h.logger.Error("get task activity failed", "error", err, "task_id", taskID)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve activities")
		return
	}

	if len(activities) == 0 {
		middleware.WriteJSON(w, http.StatusOK, map[string]any{"activities": []model.ActivityEntry{}})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, map[string]any{"activities": activities})
}

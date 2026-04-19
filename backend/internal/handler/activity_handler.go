package handler

import (
	"log"
	"net/http"

	"github.com/floqast/task-management/backend/internal/middleware"
	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/service"
)

type ActivityHandler struct {
	activityService *service.ActivityService
	logger          *log.Logger
}

func NewActivityHandler(activityService *service.ActivityService, logger *log.Logger) *ActivityHandler {
	return &ActivityHandler{
		activityService: activityService,
		logger:          logger,
	}
}

func (h *ActivityHandler) GetProjectActivity(w http.ResponseWriter, r *http.Request) {
	projectID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: readIDParam: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid project id")
	}

	activities, err := h.activityService.FindByProject(r.Context(), projectID)
	if err != nil {
		h.logger.Printf("ERROR: findByProject: %v", err)
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
		h.logger.Printf("ERROR: readIDParam: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	activities, err := h.activityService.FindByTask(r.Context(), taskID)
	if err != nil {
		h.logger.Printf("ERROR: findByTask: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve activities")
		return
	}

	if len(activities) == 0 {
		middleware.WriteJSON(w, http.StatusOK, map[string]any{"activities": []model.ActivityEntry{}})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, map[string]any{"activities": activities})
}

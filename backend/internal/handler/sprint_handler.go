package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/floqast/task-management/backend/internal/middleware"
	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SprintHandler struct {
	sprintService *service.SprintService
	taskService   *service.TaskService
	logger        *slog.Logger
}

func NewSprintHandler(sprintService *service.SprintService, taskService *service.TaskService, logger *slog.Logger) *SprintHandler {
	return &SprintHandler{
		sprintService: sprintService,
		taskService:   taskService,
		logger:        logger,
	}
}

func (h *SprintHandler) ListSprints(w http.ResponseWriter, r *http.Request) {
	projectID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid project id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid project id")
		return
	}

	activeOnly := r.URL.Query().Get("active") == "true"

	sprints, err := h.sprintService.FindByProject(r.Context(), projectID, activeOnly)
	if err != nil {
		h.logger.Error("list sprints", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve sprints")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"sprints": sprints})
}

func (h *SprintHandler) CreateSprint(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in")
		return
	}

	projectID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid project id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid project id")
		return
	}

	var req model.CreateSprintRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("decode create sprint request", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}
	if req.Name == "" {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "sprint name is required")
		return
	}
	if req.StartDate.IsZero() || req.EndDate.IsZero() {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "start_date and end_date are required")
		return
	}
	if req.EndDate.Before(req.StartDate) {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "end_date must be after start_date")
		return
	}

	sprint := &model.Sprint{
		ID:        primitive.NewObjectID().Hex(),
		ProjectID: projectID,
		Name:      req.Name,
		Label:     req.Label,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		IsActive:  true,
		TaskCount: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = h.sprintService.Create(r.Context(), sprint)
	if err != nil {
		h.logger.Error("create sprint failed", "error", err, "project_id", projectID, "user_id", user.ID)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not create sprint")
		return
	}

	h.logger.Info("sprint created", "sprint_id", sprint.ID, "name", sprint.Name, "project_id", projectID, "user_id", user.ID, "user_name", user.Name)
	middleware.WriteJSON(w, http.StatusCreated, map[string]any{"sprint": sprint})
}

func (h *SprintHandler) GetSprint(w http.ResponseWriter, r *http.Request) {
	sprintID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid sprint id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid sprint id")
		return
	}

	sprint, err := h.sprintService.FindByID(r.Context(), sprintID)
	if err != nil {
		h.logger.Error("find sprint", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve sprint")
		return
	}
	if sprint == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "sprint not found")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"sprint": sprint})
}

func (h *SprintHandler) UpdateSprint(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in")
		return
	}

	sprintID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid sprint id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid sprint id")
		return
	}

	sprint, err := h.sprintService.FindByID(r.Context(), sprintID)
	if err != nil {
		h.logger.Error("find sprint", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve sprint")
		return
	}
	if sprint == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "sprint not found")
		return
	}

	var req model.UpdateSprintRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("decode update sprint request", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}

	if req.Name != nil {
		sprint.Name = *req.Name
	}
	if req.Label != nil {
		sprint.Label = *req.Label
	}
	if req.StartDate != nil {
		sprint.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		sprint.EndDate = *req.EndDate
	}
	if req.IsActive != nil {
		sprint.IsActive = *req.IsActive
	}
	sprint.UpdatedAt = time.Now()

	err = h.sprintService.Update(r.Context(), sprint)
	if err != nil {
		h.logger.Error("update sprint failed", "error", err, "sprint_id", sprintID, "user_id", user.ID)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not update sprint")
		return
	}

	h.logger.Info("sprint updated", "sprint_id", sprintID, "user_id", user.ID, "user_name", user.Name)
	middleware.WriteJSON(w, http.StatusOK, map[string]any{"sprint": sprint})
}

func (h *SprintHandler) DeleteSprint(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in")
		return
	}

	sprintID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid sprint id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid sprint id")
		return
	}

	sprint, err := h.sprintService.FindByID(r.Context(), sprintID)
	if err != nil {
		h.logger.Error("find sprint", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve sprint")
		return
	}
	if sprint == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "sprint not found")
		return
	}

	err = h.sprintService.Delete(r.Context(), sprintID)
	if err != nil {
		h.logger.Error("delete sprint failed", "error", err, "sprint_id", sprintID, "user_id", user.ID)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not delete sprint")
		return
	}

	h.logger.Info("sprint deleted", "sprint_id", sprintID, "user_id", user.ID, "user_name", user.Name)
	middleware.WriteJSON(w, http.StatusNoContent, map[string]any{})
}

func (h *SprintHandler) AddTaskToSprint(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in")
		return
	}

	sprintID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid sprint id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid sprint id")
		return
	}

	sprint, err := h.sprintService.FindByID(r.Context(), sprintID)
	if err != nil {
		h.logger.Error("find sprint", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve sprint")
		return
	}
	if sprint == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "sprint not found")
		return
	}

	var req model.SprintTaskRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("decode sprint task request", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}
	if req.TaskID == "" {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "task_id is required")
		return
	}

	// verify task exists
	task, err := h.taskService.FindByID(r.Context(), req.TaskID)
	if err != nil {
		h.logger.Error("find task", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve task")
		return
	}
	if task == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "task not found")
		return
	}

	// update task's sprint_id
	task.SprintID = &sprintID
	task.UpdatedAt = time.Now()
	err = h.taskService.Update(r.Context(), task)
	if err != nil {
		h.logger.Error("update task sprint", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not update task")
		return
	}

	// increment sprint task count
	err = h.sprintService.AddTask(r.Context(), sprintID, req.TaskID)
	if err != nil {
		h.logger.Error("add task to sprint", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not add task to sprint")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"message": "task added to sprint"})
}

func (h *SprintHandler) RemoveTaskFromSprint(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in")
		return
	}

	sprintID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid sprint id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid sprint id")
		return
	}

	sprint, err := h.sprintService.FindByID(r.Context(), sprintID)
	if err != nil {
		h.logger.Error("find sprint", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve sprint")
		return
	}
	if sprint == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "sprint not found")
		return
	}

	var req model.SprintTaskRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("decode sprint task request", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}
	if req.TaskID == "" {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "task_id is required")
		return
	}

	// verify task exists and clear its sprint_id
	task, err := h.taskService.FindByID(r.Context(), req.TaskID)
	if err != nil {
		h.logger.Error("find task", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve task")
		return
	}
	if task == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "task not found")
		return
	}

	task.SprintID = nil
	task.UpdatedAt = time.Now()
	err = h.taskService.Update(r.Context(), task)
	if err != nil {
		h.logger.Error("update task sprint", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not update task")
		return
	}

	// decrement sprint task count
	err = h.sprintService.RemoveTask(r.Context(), sprintID, req.TaskID)
	if err != nil {
		h.logger.Error("remove task from sprint", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not remove task from sprint")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"message": "task removed from sprint"})
}

package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/floqast/task-management/backend/internal/middleware"
	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/model/dto"
	"github.com/floqast/task-management/backend/internal/repository"
	"github.com/floqast/task-management/backend/internal/service"
)

type TaskHandler struct {
	taskService    *service.TaskService
	projectService *service.ProjectService
	logger         *slog.Logger
}

func NewTaskHandler(taskService *service.TaskService, projectService *service.ProjectService, logger *slog.Logger) *TaskHandler {
	return &TaskHandler{
		taskService:    taskService,
		projectService: projectService,
		logger:         logger,
	}
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	projectID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid project id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid project id")
		return
	}

	page, pageSize := middleware.GetPaginationParams(r)
	filters := parseTaskFilters(r)

	tasks, total, err := h.taskService.FindByProject(r.Context(), projectID, filters, page, pageSize)
	if err != nil {
		h.logger.Error("list tasks failed", "error", err, "project_id", projectID)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve tasks")
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	middleware.WriteJSON(w, http.StatusOK, dto.PaginatedResponse[model.Task]{
		Data: tasks,
		Pagination: dto.Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	})
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in")
		return
	}

	projectID, err := middleware.ReadIDParam(r)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid project id")
		return
	}

	project, err := h.projectService.FindByID(r.Context(), projectID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve project")
		return
	}
	if project == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "project not found")
		return
	}

	var req dto.CreateTaskRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}
	if req.Title == "" {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "task title is required")
		return
	}

	task, err := h.taskService.CreateTask(r.Context(), user, projectID, req)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not create task")
		return
	}

	h.logger.Info("task created", "task_id", task.ID, "project_id", projectID, "user_id", user.ID)
	middleware.WriteJSON(w, http.StatusCreated, map[string]any{"task": task})
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid task id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	task, err := h.taskService.FindByID(r.Context(), taskID)
	if err != nil {
		h.logger.Error("find task", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve task")
		return
	}
	if task == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "task not found")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"task": task})
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in")
		return
	}

	taskID, err := middleware.ReadIDParam(r)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	var req dto.UpdateTaskRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}

	task, err := h.taskService.UpdateTask(r.Context(), user, taskID, req)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not update task")
		return
	}
	if task == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "task not found")
		return
	}

	h.logger.Info("task updated", "task_id", taskID, "user_id", user.ID)
	middleware.WriteJSON(w, http.StatusOK, map[string]any{"task": task})
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in")
		return
	}

	taskID, err := middleware.ReadIDParam(r)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	err = h.taskService.DeleteTask(r.Context(), user, taskID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not delete task")
		return
	}

	h.logger.Info("task deleted", "task_id", taskID, "user_id", user.ID)
	middleware.WriteJSON(w, http.StatusNoContent, map[string]any{})
}

func (h *TaskHandler) AssignTask(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in")
		return
	}

	taskID, err := middleware.ReadIDParam(r)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	var req dto.AssignTaskRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}

	task, err := h.taskService.AssignTask(r.Context(), user, taskID, req)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not assign task")
		return
	}
	if task == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "task not found")
		return
	}

	h.logger.Info("task assigned", "task_id", taskID, "user_id", user.ID)
	middleware.WriteJSON(w, http.StatusOK, map[string]any{"task": task})
}

func (h *TaskHandler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in")
		return
	}

	taskID, err := middleware.ReadIDParam(r)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	var req dto.UpdateStatusRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}

	task, err := h.taskService.UpdateStatus(r.Context(), user, taskID, req)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not update task status")
		return
	}
	if task == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "task not found")
		return
	}

	h.logger.Info("task status changed", "task_id", taskID, "user_id", user.ID)
	middleware.WriteJSON(w, http.StatusOK, map[string]any{"task": task})
}

func (h *TaskHandler) GetTaskTimeTracking(w http.ResponseWriter, r *http.Request) {
	taskID, err := middleware.ReadIDParam(r)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	task, err := h.taskService.FindByID(r.Context(), taskID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve task")
		return
	}
	if task == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "task not found")
		return
	}

	tracking := task.TimeTracking
	if tracking == nil {
		tracking = &model.TimeTracking{}
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"time_tracking": tracking})
}

func (h *TaskHandler) LogTaskTime(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in")
		return
	}

	taskID, err := middleware.ReadIDParam(r)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	var req dto.LogTimeRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}
	if req.Hours <= 0 {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "hours must be greater than 0")
		return
	}

	task, err := h.taskService.LogTime(r.Context(), user, taskID, req)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not log time")
		return
	}
	if task == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "task not found")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"time_tracking": task.TimeTracking})
}

func (h *TaskHandler) GetMyTasks(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in")
		return
	}

	page, pageSize := middleware.GetPaginationParams(r)
	filters := parseTaskFilters(r)

	tasks, total, err := h.taskService.FindByAssignee(r.Context(), user.ID, filters, page, pageSize)
	if err != nil {
		h.logger.Error("get my tasks", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve tasks")
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	middleware.WriteJSON(w, http.StatusOK, dto.PaginatedResponse[model.Task]{
		Data: tasks,
		Pagination: dto.Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	})
}

// parseTaskFilters reads optional query params into a TaskFilters struct
func parseTaskFilters(r *http.Request) repository.TaskFilters {
	var filters repository.TaskFilters

	if s := r.URL.Query().Get("status"); s != "" {
		status := model.TaskStatus(s)
		filters.Status = &status
	}
	if p := r.URL.Query().Get("priority"); p != "" {
		priority := model.Priority(p)
		filters.Priority = &priority
	}
	if a := r.URL.Query().Get("assignee_id"); a != "" {
		filters.AssigneeID = &a
	}
	if rp := r.URL.Query().Get("reporter_id"); rp != "" {
		filters.ReporterID = &rp
	}
	if sp := r.URL.Query().Get("sprint_id"); sp != "" {
		filters.SprintID = &sp
	}
	if sb := r.URL.Query().Get("sort_by"); sb != "" {
		filters.SortBy = sb
	}
	if so := r.URL.Query().Get("sort_order"); so != "" {
		filters.SortOrder = so
	}

	return filters
}

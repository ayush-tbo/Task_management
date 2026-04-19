package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/floqast/task-management/backend/internal/middleware"
	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
	"github.com/floqast/task-management/backend/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskHandler struct {
	taskService    *service.TaskService
	projectService *service.ProjectService
	logger         *log.Logger
}

func NewTaskHandler(taskService *service.TaskService, projectService *service.ProjectService, logger *log.Logger) *TaskHandler {
	return &TaskHandler{
		taskService:    taskService,
		projectService: projectService,
		logger:         logger,
	}
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	projectID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: invalid project id: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid project id")
		return
	}

	page, pageSize := middleware.GetPaginationParams(r)
	filters := parseTaskFilters(r)

	tasks, total, err := h.taskService.FindByProject(r.Context(), projectID, filters, page, pageSize)
	if err != nil {
		h.logger.Printf("ERROR: list tasks: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve tasks")
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	middleware.WriteJSON(w, http.StatusOK, model.PaginatedResponse[model.Task]{
		Data: tasks,
		Pagination: model.Pagination{
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
		h.logger.Printf("ERROR: invalid project id: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid project id")
		return
	}

	project, err := h.projectService.FindByID(r.Context(), projectID)
	if err != nil {
		h.logger.Printf("ERROR: find project: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve project")
		return
	}
	if project == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "project not found")
		return
	}

	var req model.CreateTaskRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decode create task request: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}
	if req.Title == "" {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "task title is required")
		return
	}

	status := req.Status
	if status == "" {
		status = model.StatusTodo
	}
	priority := req.Priority
	if priority == "" {
		priority = model.PriorityP3
	}

	task := &model.Task{
		ID:          primitive.NewObjectID().Hex(),
		ProjectID:   projectID,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		Priority:    priority,
		DueDate:     req.DueDate,
		AssigneeID:  req.AssigneeID,
		ReporterID:  user.ID,
		LabelIDs:    req.LabelIDs,
		SprintID:    req.SprintID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = h.taskService.Create(r.Context(), task)
	if err != nil {
		h.logger.Printf("ERROR: create task: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not create task")
		return
	}

	middleware.WriteJSON(w, http.StatusCreated, map[string]any{"task": task})
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: invalid task id: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	task, err := h.taskService.FindByID(r.Context(), taskID)
	if err != nil {
		h.logger.Printf("ERROR: find task: %v", err)
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
		h.logger.Printf("ERROR: invalid task id: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	task, err := h.taskService.FindByID(r.Context(), taskID)
	if err != nil {
		h.logger.Printf("ERROR: find task: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve task")
		return
	}
	if task == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "task not found")
		return
	}

	var req model.UpdateTaskRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decode update task request: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if req.AssigneeID != nil {
		task.AssigneeID = req.AssigneeID
	}
	if req.LabelIDs != nil {
		task.LabelIDs = req.LabelIDs
	}
	if req.SprintID != nil {
		task.SprintID = req.SprintID
	}
	task.UpdatedAt = time.Now()

	err = h.taskService.Update(r.Context(), task)
	if err != nil {
		h.logger.Printf("ERROR: update task: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not update task")
		return
	}

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
		h.logger.Printf("ERROR: invalid task id: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	task, err := h.taskService.FindByID(r.Context(), taskID)
	if err != nil {
		h.logger.Printf("ERROR: find task: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve task")
		return
	}
	if task == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "task not found")
		return
	}

	err = h.taskService.Delete(r.Context(), taskID)
	if err != nil {
		h.logger.Printf("ERROR: delete task: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not delete task")
		return
	}

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
		h.logger.Printf("ERROR: invalid task id: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	task, err := h.taskService.FindByID(r.Context(), taskID)
	if err != nil {
		h.logger.Printf("ERROR: find task: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve task")
		return
	}
	if task == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "task not found")
		return
	}

	var req model.AssignTaskRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decode assign task request: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}
	if req.AssigneeID == "" {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "assignee_id is required")
		return
	}

	task.AssigneeID = &req.AssigneeID
	task.UpdatedAt = time.Now()

	err = h.taskService.Update(r.Context(), task)
	if err != nil {
		h.logger.Printf("ERROR: assign task: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not assign task")
		return
	}

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
		h.logger.Printf("ERROR: invalid task id: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	task, err := h.taskService.FindByID(r.Context(), taskID)
	if err != nil {
		h.logger.Printf("ERROR: find task: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve task")
		return
	}
	if task == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "task not found")
		return
	}

	var req model.UpdateStatusRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decode update status request: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}

	task.Status = req.Status
	task.UpdatedAt = time.Now()

	err = h.taskService.Update(r.Context(), task)
	if err != nil {
		h.logger.Printf("ERROR: update task status: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not update task status")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"task": task})
}

func (h *TaskHandler) GetTaskTimeTracking(w http.ResponseWriter, r *http.Request) {
	taskID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: invalid task id: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	task, err := h.taskService.FindByID(r.Context(), taskID)
	if err != nil {
		h.logger.Printf("ERROR: find task: %v", err)
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
		h.logger.Printf("ERROR: invalid task id: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	task, err := h.taskService.FindByID(r.Context(), taskID)
	if err != nil {
		h.logger.Printf("ERROR: find task: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve task")
		return
	}
	if task == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "task not found")
		return
	}

	var req model.LogTimeRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decode log time request: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}
	if req.Hours <= 0 {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "hours must be greater than 0")
		return
	}

	if task.TimeTracking == nil {
		task.TimeTracking = &model.TimeTracking{}
	}
	task.TimeTracking.LoggedHours += req.Hours
	task.UpdatedAt = time.Now()

	err = h.taskService.Update(r.Context(), task)
	if err != nil {
		h.logger.Printf("ERROR: log time: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not log time")
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
		h.logger.Printf("ERROR: get my tasks: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve tasks")
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	middleware.WriteJSON(w, http.StatusOK, model.PaginatedResponse[model.Task]{
		Data: tasks,
		Pagination: model.Pagination{
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

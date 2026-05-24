package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/floqast/task-management/backend/internal/middleware"
	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/model/dto"
	"github.com/floqast/task-management/backend/internal/service"
)

type ProjectHandler struct {
	projectService *service.ProjectService
	taskService    *service.TaskService
	logger         *slog.Logger
}

func NewProjectHandler(projectService *service.ProjectService, taskService *service.TaskService, logger *slog.Logger) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		taskService:    taskService,
		logger:         logger,
	}
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in to create a project")
		return
	}
	page, pageSize := middleware.GetPaginationParams(r)
	projects, total, err := h.projectService.FindByUser(r.Context(), user.ID, page, pageSize)
	if err != nil {
		h.logger.Error("list projects failed", "error", err, "user_id", user.ID)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve projects")
		return
	}
	totalPages := (total + pageSize - 1) / pageSize
	middleware.WriteJSON(w, http.StatusOK, dto.PaginatedResponse[model.Project]{
		Data: projects,
		Pagination: dto.Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	})

}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in to create a project")
		return
	}

	var req dto.CreateProjectRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}
	if req.Name == "" {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "project name is required")
		return
	}

	project, err := h.projectService.CreateProject(r.Context(), user, req)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not create project")
		return
	}

	h.logger.Info("project created", "project_id", project.ID, "user_id", user.ID)
	middleware.WriteJSON(w, http.StatusCreated, map[string]any{"project": project})
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in to create a project")
		return
	}
	projectID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid project id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid project id")
		return
	}
	if projectID == "" {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "project id is required")
		return
	}
	project, err := h.projectService.FindByID(r.Context(), projectID)
	if err != nil {
		h.logger.Error("find project by id", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve project")
		return
	}
	if project == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "project not found")
		return
	}
	middleware.WriteJSON(w, http.StatusOK, map[string]any{"project": project})
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
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

	project, err := h.projectService.FindByID(r.Context(), projectID)
	if err != nil {
		h.logger.Error("find project", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve project")
		return
	}
	if project == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "project not found")
		return
	}

	if project.OwnerID != user.ID {
		middleware.WriteError(w, http.StatusForbidden, "forbidden", "you are not the owner of this project")
		return
	}

	var req dto.UpdateProjectRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("decode update project request", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}

	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.Description != nil {
		project.Description = *req.Description
	}
	project.UpdatedAt = time.Now()

	err = h.projectService.Update(r.Context(), project)
	if err != nil {
		h.logger.Error("update project failed", "error", err, "project_id", projectID, "user_id", user.ID)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not update project")
		return
	}

	h.logger.Info("project updated", "project_id", projectID, "user_id", user.ID, "user_name", user.Name)
	middleware.WriteJSON(w, http.StatusOK, map[string]any{"project": project})
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
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

	project, err := h.projectService.FindByID(r.Context(), projectID)
	if err != nil {
		h.logger.Error("find project", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve project")
		return
	}
	if project == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "project not found")
		return
	}

	if project.OwnerID != user.ID {
		middleware.WriteError(w, http.StatusForbidden, "forbidden", "you are not the owner of this project")
		return
	}

	err = h.projectService.Delete(r.Context(), projectID)
	if err != nil {
		h.logger.Error("delete project failed", "error", err, "project_id", projectID, "user_id", user.ID)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not delete project")
		return
	}

	h.logger.Info("project deleted", "project_id", projectID, "user_id", user.ID, "user_name", user.Name)
	middleware.WriteJSON(w, http.StatusNoContent, map[string]any{})
}

func (h *ProjectHandler) ListProjectMembers(w http.ResponseWriter, r *http.Request) {
	projectID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid project id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid project id")
		return
	}

	members, err := h.projectService.ListMembers(r.Context(), projectID)
	if err != nil {
		h.logger.Error("list members", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve members")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"members": members})
}

func (h *ProjectHandler) AddProjectMember(w http.ResponseWriter, r *http.Request) {
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

	var req dto.AddMemberRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}
	if req.UserID == "" {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "user_id is required")
		return
	}

	member, err := h.projectService.AddMember(r.Context(), user, projectID, req)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not add member")
		return
	}
	if member == nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "user not found")
		return
	}

	h.logger.Info("member added", "project_id", projectID, "user_id", user.ID)
	middleware.WriteJSON(w, http.StatusCreated, map[string]any{"member": member})
}

func (h *ProjectHandler) RemoveProjectMember(w http.ResponseWriter, r *http.Request) {
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

	userID, err := middleware.ReadURLParam(r, "userId")
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid user id")
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
	if project.OwnerID != user.ID {
		middleware.WriteError(w, http.StatusForbidden, "forbidden", "only the owner can remove members")
		return
	}

	err = h.projectService.RemoveMember(r.Context(), user, projectID, userID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not remove member")
		return
	}

	h.logger.Info("member removed", "project_id", projectID, "user_id", user.ID)
	middleware.WriteJSON(w, http.StatusNoContent, map[string]any{})
}

func (h *ProjectHandler) GetStatusChart(w http.ResponseWriter, r *http.Request) {
	projectID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid project id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid project id")
		return
	}

	entries, err := h.taskService.CountByStatus(r.Context(), projectID)
	if err != nil {
		h.logger.Error("count by status", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve chart")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, dto.StatusChart{
		ProjectID: projectID,
		Data:      entries,
	})
}

func (h *ProjectHandler) GetPriorityChart(w http.ResponseWriter, r *http.Request) {
	projectID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid project id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid project id")
		return
	}

	entries, err := h.taskService.CountByPriority(r.Context(), projectID)
	if err != nil {
		h.logger.Error("count by priority", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve chart")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, dto.PriorityChart{
		ProjectID: projectID,
		Data:      entries,
	})
}

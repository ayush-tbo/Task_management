package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/floqast/task-management/backend/internal/middleware"
	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectHandler struct {
	projectService      *service.ProjectService
	taskService         *service.TaskService
	activityService     *service.ActivityService
	userService         *service.UserService
	notificationService *service.NotificationService
	logger              *slog.Logger
}

func NewProjectHandler(projectService *service.ProjectService, taskService *service.TaskService, activityService *service.ActivityService, userService *service.UserService, notificationService *service.NotificationService, logger *slog.Logger) *ProjectHandler {
	return &ProjectHandler{
		projectService:      projectService,
		taskService:         taskService,
		activityService:     activityService,
		userService:         userService,
		notificationService: notificationService,
		logger:              logger,
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
	middleware.WriteJSON(w, http.StatusOK, model.PaginatedResponse[model.Project]{
		Data: projects,
		Pagination: model.Pagination{
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
	var req model.CreateProjectRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("decode create project request", "error", err, "user_id", user.ID)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}
	if req.Name == "" {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "project name is required")
		return
	}
	project := &model.Project{
		ID:          primitive.NewObjectID().Hex(),
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     user.ID,
		MemberCount: 1,
		TaskCount:   0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = h.projectService.Create(r.Context(), project)
	if err != nil {
		h.logger.Error("create project failed", "error", err, "user_id", user.ID, "project_name", req.Name)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not create project")
		return
	}

	// add creator as owner member
	owner := &model.ProjectMember{
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		Role:      model.RoleOwner,
		JoinedAt:  time.Now(),
	}
	err = h.projectService.AddMember(r.Context(), project.ID, owner)
	if err != nil {
		h.logger.Error("add owner member failed", "error", err, "user_id", user.ID, "project_id", project.ID)
	}

	h.logger.Info("project created", "project_id", project.ID, "project_name", project.Name, "user_id", user.ID, "user_name", user.Name)

	go func() {
		entry := &model.ActivityEntry{
			ID:        primitive.NewObjectID().Hex(),
			ProjectID: project.ID,
			UserID:    user.ID,
			User:      user,
			Action:    model.ActionMemberAdded,
			Details:   map[string]interface{}{"project_name": project.Name},
			CreatedAt: time.Now(),
		}
		if err := h.activityService.Create(context.Background(), entry); err != nil {
			h.logger.Error("activity log failed", "error", err, "project_id", project.ID)
		}
	}()

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

	var req model.UpdateProjectRequest
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

	var req model.AddMemberRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("decode add member request", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}
	if req.UserID == "" {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "user_id is required")
		return
	}

	role := req.Role
	if role == "" {
		role = model.RoleMember
	}

	// Look up the user to populate name/email/avatar
	userInfo, err := h.userService.FindByID(r.Context(), req.UserID)
	if err != nil || userInfo == nil {
		h.logger.Error("find user for member", "error", err, "user_id", req.UserID)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "user not found")
		return
	}

	member := &model.ProjectMember{
		UserID:    req.UserID,
		Name:      userInfo.Name,
		Email:     userInfo.Email,
		AvatarURL: userInfo.AvatarURL,
		Role:      role,
		JoinedAt:  time.Now(),
	}

	err = h.projectService.AddMember(r.Context(), projectID, member)
	if err != nil {
		h.logger.Error("add member", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not add member")
		return
	}

	_ = h.projectService.IncrementMemberCount(r.Context(), projectID, 1)

	go func() {
		// Notify the added member
		if req.UserID != user.ID {
			refType := "project"
			n := &model.Notification{
				ID:            primitive.NewObjectID().Hex(),
				UserID:        req.UserID,
				Type:          model.NotifAlert,
				Title:         "Added to Project",
				Message:       user.Name + " added you to project \"" + project.Name + "\"",
				ReferenceType: &refType,
				ReferenceID:   &projectID,
				CreatedAt:     time.Now(),
			}
			if err := h.notificationService.Create(context.Background(), n); err != nil {
				h.logger.Error("notification create failed", "error", err)
			}
		}

		entry := &model.ActivityEntry{
			ID:        primitive.NewObjectID().Hex(),
			ProjectID: projectID,
			UserID:    user.ID,
			User:      user,
			Action:    model.ActionMemberAdded,
			Details:   map[string]interface{}{"member_user_id": req.UserID, "role": string(role)},
			CreatedAt: time.Now(),
		}
		if err := h.activityService.Create(context.Background(), entry); err != nil {
			h.logger.Error("activity log", "error", err)
		}
	}()

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
		h.logger.Error("invalid project id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid project id")
		return
	}

	userID, err := middleware.ReadURLParam(r, "userId")
	if err != nil {
		h.logger.Error("invalid user id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid user id")
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
		middleware.WriteError(w, http.StatusForbidden, "forbidden", "only the owner can remove members")
		return
	}

	err = h.projectService.RemoveMember(r.Context(), projectID, userID)
	if err != nil {
		h.logger.Error("remove member", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not remove member")
		return
	}

	_ = h.projectService.IncrementMemberCount(r.Context(), projectID, -1)

	go func() {
		entry := &model.ActivityEntry{
			ID:        primitive.NewObjectID().Hex(),
			ProjectID: projectID,
			UserID:    user.ID,
			User:      user,
			Action:    model.ActionMemberRemoved,
			Details:   map[string]interface{}{"removed_user_id": userID},
			CreatedAt: time.Now(),
		}
		if err := h.activityService.Create(context.Background(), entry); err != nil {
			h.logger.Error("activity log", "error", err)
		}
	}()

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

	middleware.WriteJSON(w, http.StatusOK, model.StatusChart{
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

	middleware.WriteJSON(w, http.StatusOK, model.PriorityChart{
		ProjectID: projectID,
		Data:      entries,
	})
}

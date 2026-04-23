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

type CommentHandler struct {
	commentService      *service.CommentService
	taskService         *service.TaskService
	activity            *service.ActivityService
	notificationService *service.NotificationService
	logger              *slog.Logger
}

func NewCommentHandler(commentService *service.CommentService, taskService *service.TaskService, activity *service.ActivityService, notificationService *service.NotificationService, logger *slog.Logger) *CommentHandler {
	return &CommentHandler{
		commentService:      commentService,
		taskService:         taskService,
		activity:            activity,
		notificationService: notificationService,
		logger:              logger,
	}
}

func (h *CommentHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	taskID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("readIDParam", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	comments, err := h.commentService.FindByTask(r.Context(), taskID)
	if err != nil {
		h.logger.Error("findByTask", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve comments")
		return
	}

	if len(comments) == 0 {
		middleware.WriteJSON(w, http.StatusOK, map[string]any{"comments": []model.Comment{}})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, map[string]any{"comments": comments})
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	taskID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("readIDParam", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in to comment")
		return
	}

	var req model.CreateCommentRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("decoding create comment request", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request payload")
		return
	}

	comment := &model.Comment{
		ID:        primitive.NewObjectID().Hex(),
		TaskID:    taskID,
		UserID:    user.ID,
		User:      user,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = h.commentService.Create(r.Context(), comment)
	if err != nil {
		h.logger.Error("create comment failed", "error", err, "task_id", taskID, "user_id", user.ID)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	h.logger.Info("comment created", "comment_id", comment.ID, "task_id", taskID, "user_id", user.ID, "user_name", user.Name)

	// Notify task reporter and assignee about the new comment
	go func() {
		task, err := h.taskService.FindByID(context.Background(), taskID)
		if err == nil && task != nil {
			refType := "task"
			notified := map[string]bool{user.ID: true}
			notify := func(targetUserID string) {
				if notified[targetUserID] {
					return
				}
				notified[targetUserID] = true
				n := &model.Notification{
					ID:            primitive.NewObjectID().Hex(),
					UserID:        targetUserID,
					Type:          model.NotifMention,
					Title:         "New Comment",
					Message:       user.Name + " commented on \"" + task.Title + "\"",
					ReferenceType: &refType,
					ReferenceID:   &task.ID,
					CreatedAt:     time.Now(),
				}
				if err := h.notificationService.Create(context.Background(), n); err != nil {
					h.logger.Error("notification create failed", "error", err)
				}
			}
			if task.ReporterID != "" {
				notify(task.ReporterID)
			}
			if task.AssigneeID != nil && *task.AssigneeID != "" {
				notify(*task.AssigneeID)
			}
		}
	}()

	activity := &model.ActivityEntry{
		ID:        primitive.NewObjectID().Hex(),
		ProjectID: req.ProjectID,
		TaskID:    &taskID,
		UserID:    user.ID,
		User:      user,
		Action:    model.ActionCommentAdded,
		Details:   map[string]interface{}{"comment_id": comment.ID},
		CreatedAt: time.Now(),
	}

	go func() {
		err := h.activity.Create(context.Background(), activity)

		if err != nil {
			h.logger.Error("Activity Logging Error", "error", err)
		}
	}()

	middleware.WriteJSON(w, http.StatusCreated, map[string]any{"comment": comment})
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	commentID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("readIDParam", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid comment id")
		return
	}

	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in to update comment")
		return
	}

	existingComment, err := h.commentService.FindByID(r.Context(), commentID)
	if err != nil {
		h.logger.Error("findByID", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	if existingComment == nil {
		http.NotFound(w, r)
		return
	}

	if existingComment.UserID != user.ID {
		middleware.WriteError(w, http.StatusForbidden, "Forbidden", "You cannot update this comment")
		return
	}

	oldContent := existingComment.Content

	var input model.UpdateCommentRequest

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.logger.Error("decoding update comment request", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request payload")
		return
	}

	if input.Content != nil {
		existingComment.Content = *input.Content
	}
	existingComment.UpdatedAt = time.Now()

	err = h.commentService.Update(r.Context(), existingComment)
	if err != nil {
		h.logger.Error("update comment failed", "error", err, "comment_id", commentID, "user_id", user.ID)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	h.logger.Info("comment updated", "comment_id", commentID, "user_id", user.ID, "user_name", user.Name)

	activity := &model.ActivityEntry{
		ID:        primitive.NewObjectID().Hex(),
		ProjectID: input.ProjectID,
		TaskID:    &input.TaskID,
		UserID:    user.ID,
		User:      user,
		Action:    model.ActionCommentChanged,
		Details:   map[string]interface{}{"comment_id": existingComment.ID, "old_content": oldContent},
		CreatedAt: time.Now(),
	}

	go func() {
		err := h.activity.Create(context.Background(), activity)

		if err != nil {
			h.logger.Error("Activity Logging Error", "error", err)
		}
	}()

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"comment": existingComment})
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("readIDParam", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid comment id")
		return
	}

	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in to delete comment")
		return
	}

	existingComment, err := h.commentService.FindByID(r.Context(), commentID)
	if err != nil {
		h.logger.Error("findByID", "error", err)
		middleware.WriteError(w, http.StatusNotFound, "notFound", "Comment not found")
		return
	}

	if existingComment.UserID != user.ID {
		middleware.WriteError(w, http.StatusForbidden, "Forbidden", "You cannot delete this comment")
		return
	}

	var input model.DeleteCommentRequest

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.logger.Error("decoding delete comment request", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request payload")
		return
	}

	err = h.commentService.Delete(r.Context(), commentID)
	if err != nil {
		h.logger.Error("delete comment failed", "error", err, "comment_id", commentID, "user_id", user.ID)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "delete failed")
	}

	h.logger.Info("comment deleted", "comment_id", commentID, "user_id", user.ID, "user_name", user.Name)

	activity := &model.ActivityEntry{
		ID:        primitive.NewObjectID().Hex(),
		ProjectID: input.ProjectID,
		TaskID:    &input.TaskID,
		UserID:    user.ID,
		User:      user,
		Action:    model.ActionCommentDeleted,
		Details:   map[string]interface{}{"comment_id": existingComment.ID, "info": "Comment was removed by the author"},
		CreatedAt: time.Now(),
	}

	go func() {
		err := h.activity.Create(context.Background(), activity)

		if err != nil {
			h.logger.Error("Activity Logging Error", "error", err)
		}
	}()

	middleware.WriteJSON(w, http.StatusNoContent, map[string]any{})
}

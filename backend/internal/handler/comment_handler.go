package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/floqast/task-management/backend/internal/middleware"
	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/model/dto"
	"github.com/floqast/task-management/backend/internal/service"
)

type CommentHandler struct {
	commentService *service.CommentService
	logger         *slog.Logger
}

func NewCommentHandler(commentService *service.CommentService, logger *slog.Logger) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		logger:         logger,
	}
}

func (h *CommentHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	taskID, err := middleware.ReadIDParam(r)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	comments, err := h.commentService.FindByTask(r.Context(), taskID)
	if err != nil {
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
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid task id")
		return
	}

	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in to comment")
		return
	}

	var req dto.CreateCommentRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request payload")
		return
	}

	comment, err := h.commentService.CreateComment(r.Context(), user, taskID, req)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not create comment")
		return
	}

	h.logger.Info("comment created", "comment_id", comment.ID, "task_id", taskID, "user_id", user.ID)
	middleware.WriteJSON(w, http.StatusCreated, map[string]any{"comment": comment})
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	commentID, err := middleware.ReadIDParam(r)
	if err != nil {
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

	var req dto.UpdateCommentRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request payload")
		return
	}

	updated, err := h.commentService.UpdateComment(r.Context(), user, existingComment, req)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not update comment")
		return
	}

	h.logger.Info("comment updated", "comment_id", commentID, "user_id", user.ID)
	middleware.WriteJSON(w, http.StatusOK, map[string]any{"comment": updated})
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID, err := middleware.ReadIDParam(r)
	if err != nil {
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
		middleware.WriteError(w, http.StatusNotFound, "notFound", "Comment not found")
		return
	}
	if existingComment.UserID != user.ID {
		middleware.WriteError(w, http.StatusForbidden, "Forbidden", "You cannot delete this comment")
		return
	}

	var req dto.DeleteCommentRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request payload")
		return
	}

	err = h.commentService.DeleteComment(r.Context(), user, existingComment, req)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "delete failed")
		return
	}

	h.logger.Info("comment deleted", "comment_id", commentID, "user_id", user.ID)
	middleware.WriteJSON(w, http.StatusNoContent, map[string]any{})
}

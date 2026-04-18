package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/floqast/task-management/backend/internal/middleware"
	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentHandler struct {
	commentService *service.CommentService
	logger         *log.Logger
}

func NewCommentHandler(commentService *service.CommentService, logger *log.Logger) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		logger:         logger,
	}
}

func (h *CommentHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	taskID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: readIDParam: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid user id")
		return
	}

	comments, err := h.commentService.FindByTask(r.Context(), taskID)
	if err != nil {
		h.logger.Printf("ERROR: findByTask: %v", err)
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
		h.logger.Printf("ERROR: readIDParam: %v", err)
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
		h.logger.Printf("ERROR: decoding create comment request: %v", err)
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
		h.logger.Printf("ERROR: create %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	middleware.WriteJSON(w, http.StatusCreated, map[string]any{"comment": comment})
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	commentID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: readIDParam: %v", err)
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
		h.logger.Printf("ERROR: findByID: %v", err)
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

	var input model.UpdateCommentRequest

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.logger.Printf("ERROR: decoding update comment request: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request payload")
		return
	}

	if input.Content != nil {
		existingComment.Content = *input.Content
	}
	existingComment.UpdatedAt = time.Now()

	err = h.commentService.Update(r.Context(), existingComment)
	if err != nil {
		h.logger.Printf("ERROR: update: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"comment": existingComment})
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: readIDParam: %v", err)
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
		h.logger.Printf("ERROR: findByID: %v", err)
		middleware.WriteError(w, http.StatusNotFound, "notFound", "Comment not found")
		return
	}

	if existingComment.UserID != user.ID {
		middleware.WriteError(w, http.StatusForbidden, "Forbidden", "You cannot delete this comment")
		return
	}

	err = h.commentService.Delete(r.Context(), commentID)
	if err != nil {
		h.logger.Printf("ERROR: delete: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "delete failed")
	}

	middleware.WriteJSON(w, http.StatusNoContent, map[string]any{})
}

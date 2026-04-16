package handler

import (
	"net/http"

	"github.com/floqast/task-management/backend/internal/middleware"
)

type CommentHandler struct{}

func NewCommentHandler() *CommentHandler { return &CommentHandler{} }

func (h *CommentHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

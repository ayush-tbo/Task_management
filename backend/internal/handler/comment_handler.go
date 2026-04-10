package handler

import "net/http"

type CommentHandler struct{}

func NewCommentHandler() *CommentHandler { return &CommentHandler{} }

func (h *CommentHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

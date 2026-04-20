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

type LabelHandler struct {
	labelService *service.LabelService
	logger       *slog.Logger
}

func NewLabelHandler(labelService *service.LabelService, logger *slog.Logger) *LabelHandler {
	return &LabelHandler{
		labelService: labelService,
		logger:       logger,
	}
}

func (h *LabelHandler) ListLabels(w http.ResponseWriter, r *http.Request) {
	projectID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid project id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid project id")
		return
	}

	labels, err := h.labelService.FindByProject(r.Context(), projectID)
	if err != nil {
		h.logger.Error("list labels", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve labels")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"labels": labels})
}

func (h *LabelHandler) CreateLabel(w http.ResponseWriter, r *http.Request) {
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

	var req model.CreateLabelRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("decode create label request", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}
	if req.Name == "" {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "label name is required")
		return
	}

	label := &model.Label{
		ID:        primitive.NewObjectID().Hex(),
		ProjectID: projectID,
		Name:      req.Name,
		Color:     req.Color,
		CreatedAt: time.Now(),
	}

	err = h.labelService.Create(r.Context(), label)
	if err != nil {
		h.logger.Error("create label", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not create label")
		return
	}

	middleware.WriteJSON(w, http.StatusCreated, map[string]any{"label": label})
}

func (h *LabelHandler) UpdateLabel(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in")
		return
	}

	labelID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid label id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid label id")
		return
	}

	label, err := h.labelService.FindByID(r.Context(), labelID)
	if err != nil {
		h.logger.Error("find label", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve label")
		return
	}
	if label == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "label not found")
		return
	}

	var req model.CreateLabelRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("decode update label request", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request body")
		return
	}

	if req.Name != "" {
		label.Name = req.Name
	}
	if req.Color != "" {
		label.Color = req.Color
	}

	err = h.labelService.Update(r.Context(), label)
	if err != nil {
		h.logger.Error("update label", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not update label")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"label": label})
}

func (h *LabelHandler) DeleteLabel(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in")
		return
	}

	labelID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("invalid label id", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid label id")
		return
	}

	label, err := h.labelService.FindByID(r.Context(), labelID)
	if err != nil {
		h.logger.Error("find label", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve label")
		return
	}
	if label == nil {
		middleware.WriteError(w, http.StatusNotFound, "not found", "label not found")
		return
	}

	err = h.labelService.Delete(r.Context(), labelID)
	if err != nil {
		h.logger.Error("delete label", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not delete label")
		return
	}

	middleware.WriteJSON(w, http.StatusNoContent, map[string]any{})
}

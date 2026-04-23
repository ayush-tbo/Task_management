package handler

import (
	"log/slog"
	"net/http"

	"github.com/floqast/task-management/backend/internal/middleware"
	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/service"
)

type NotificationHandler struct {
	service *service.NotificationService
	logger  *slog.Logger
}

func NewNotificationHandler(service *service.NotificationService, logger *slog.Logger) *NotificationHandler {
	return &NotificationHandler{
		service: service,
		logger:  logger,
	}
}

func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in to view notifications")
		return
	}

	notifications, err := h.service.FindByUser(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("list notifications failed", "error", err, "user_id", user.ID)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "Unable to fetch notifications from server")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"notifications": notifications})
}

func (h *NotificationHandler) MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	notificationID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("readIDParam", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid notification id")
		return
	}

	err = h.service.MarkRead(r.Context(), notificationID)
	if err != nil {
		h.logger.Error("mark notification read failed", "error", err, "notification_id", notificationID)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "server not able to perform operation")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (h *NotificationHandler) MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil || user == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in to mark all notifications read")
		return
	}

	err := h.service.MarkAllRead(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("mark all notifications read failed", "error", err, "user_id", user.ID)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "server not able to perform operation")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"success": true})
}

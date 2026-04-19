package repository

import (
	"context"

	"github.com/floqast/task-management/backend/internal/model"
)

type NotificationRepository interface {
	FindByUser(ctx context.Context, userID string, unreadOnly bool, notifType *model.NotificationType, page, pageSize int) ([]model.Notification, int, int, error)
	MarkRead(ctx context.Context, id string) error
	MarkAllRead(ctx context.Context, userID string) error
	Create(ctx context.Context, notification *model.Notification) error
}

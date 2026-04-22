package service

import (
	"context"
	"log/slog"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
)

type NotificationService struct {
	repo   repository.NotificationRepository
	logger *slog.Logger
}

func NewNotificationService(repo repository.NotificationRepository, logger *slog.Logger) *NotificationService {
	return &NotificationService{
		repo:   repo,
		logger: logger,
	}
}

func (s *NotificationService) FindByUser(ctx context.Context, userID string) ([]model.Notification, error) {
	notifications, err := s.repo.FindByUser(ctx, userID)
	if err != nil {
		s.logger.Error("service: find notifications by user", "error", err, "user_id", userID)
	}
	return notifications, err
}

func (s *NotificationService) MarkRead(ctx context.Context, id string) error {
	err := s.repo.MarkRead(ctx, id)
	if err != nil {
		s.logger.Error("service: mark notification read", "error", err, "notification_id", id)
	}
	return err
}

func (s *NotificationService) MarkAllRead(ctx context.Context, userID string) error {
	err := s.repo.MarkAllRead(ctx, userID)
	if err != nil {
		s.logger.Error("service: mark all notifications read", "error", err, "user_id", userID)
	}
	return err
}

func (s *NotificationService) Create(ctx context.Context, n *model.Notification) error {
	err := s.repo.Create(ctx, n)
	if err != nil {
		s.logger.Error("service: create notification", "error", err, "user_id", n.UserID)
	}
	return err
}

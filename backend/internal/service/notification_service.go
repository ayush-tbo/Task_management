package service

import (
	"context"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
)

type NotificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

func (s *NotificationService) FindByUser(ctx context.Context, userID string) ([]model.Notification, error) {
	return s.repo.FindByUser(ctx, userID)
}

func (s *NotificationService) MarkRead(ctx context.Context, id string) error {
	return s.repo.MarkRead(ctx, id)
}

func (s *NotificationService) MarkAllRead(ctx context.Context, userID string) error {
	return s.repo.MarkAllRead(ctx, userID)
}

func (s *NotificationService) Create(ctx context.Context, n *model.Notification) error {
	return s.repo.Create(ctx, n)
}

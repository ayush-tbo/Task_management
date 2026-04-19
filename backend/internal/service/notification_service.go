package service

import "github.com/floqast/task-management/backend/internal/repository"

type NotificationService struct {
	repo repository.NotificationRepository
}

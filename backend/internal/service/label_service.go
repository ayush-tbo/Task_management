package service

import "github.com/floqast/task-management/backend/internal/repository"

type LabelService struct {
	repo     repository.LabelRepository
	activity repository.ActivityRepository
}

package service

import (
	"context"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
)

type ActivityService struct {
	repo repository.ActivityRepository
}

func NewActivityService(repo repository.ActivityRepository) *ActivityService {
	return &ActivityService{
		repo: repo,
	}
}

func (s *ActivityService) FindByProject(ctx context.Context, projectID string) ([]model.ActivityEntry, error) {
	return s.repo.FindByProject(ctx, projectID)
}

func (s *ActivityService) FindByTask(ctx context.Context, taskID string) ([]model.ActivityEntry, error) {
	return s.repo.FindByTask(ctx, taskID)
}

func (s *ActivityService) Create(ctx context.Context, activity *model.ActivityEntry) error {
	return s.repo.Create(ctx, activity)
}

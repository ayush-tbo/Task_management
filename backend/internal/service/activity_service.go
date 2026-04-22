package service

import (
	"context"
	"log/slog"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
)

type ActivityService struct {
	repo   repository.ActivityRepository
	logger *slog.Logger
}

func NewActivityService(repo repository.ActivityRepository, logger *slog.Logger) *ActivityService {
	return &ActivityService{
		repo:   repo,
		logger: logger,
	}
}

func (s *ActivityService) FindByProject(ctx context.Context, projectID string) ([]model.ActivityEntry, error) {
	entries, err := s.repo.FindByProject(ctx, projectID)
	if err != nil {
		s.logger.Error("service: find activity by project", "error", err, "project_id", projectID)
	}
	return entries, err
}

func (s *ActivityService) FindByTask(ctx context.Context, taskID string) ([]model.ActivityEntry, error) {
	entries, err := s.repo.FindByTask(ctx, taskID)
	if err != nil {
		s.logger.Error("service: find activity by task", "error", err, "task_id", taskID)
	}
	return entries, err
}

func (s *ActivityService) Create(ctx context.Context, activity *model.ActivityEntry) error {
	err := s.repo.Create(ctx, activity)
	if err != nil {
		s.logger.Error("service: create activity", "error", err, "project_id", activity.ProjectID, "action", activity.Action)
	}
	return err
}

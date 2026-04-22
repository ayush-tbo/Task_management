package service

import (
	"context"
	"log/slog"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
)

type SprintService struct {
	repo     repository.SprintRepository
	activity repository.ActivityRepository
	logger   *slog.Logger
}

func NewSprintService(repo repository.SprintRepository, activity repository.ActivityRepository, logger *slog.Logger) *SprintService {
	return &SprintService{
		repo:     repo,
		activity: activity,
		logger:   logger,
	}
}

func (s *SprintService) FindByProject(ctx context.Context, projectID string, activeOnly bool) ([]model.Sprint, error) {
	sprints, err := s.repo.FindByProject(ctx, projectID, activeOnly)
	if err != nil {
		s.logger.Error("service: find sprints by project", "error", err, "project_id", projectID)
	}
	return sprints, err
}

func (s *SprintService) FindByID(ctx context.Context, id string) (*model.Sprint, error) {
	sprint, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("service: find sprint by id", "error", err, "sprint_id", id)
	}
	return sprint, err
}

func (s *SprintService) Create(ctx context.Context, sprint *model.Sprint) error {
	err := s.repo.Create(ctx, sprint)
	if err != nil {
		s.logger.Error("service: create sprint", "error", err, "sprint_id", sprint.ID, "project_id", sprint.ProjectID)
	}
	return err
}

func (s *SprintService) Update(ctx context.Context, sprint *model.Sprint) error {
	err := s.repo.Update(ctx, sprint)
	if err != nil {
		s.logger.Error("service: update sprint", "error", err, "sprint_id", sprint.ID)
	}
	return err
}

func (s *SprintService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("service: delete sprint", "error", err, "sprint_id", id)
	}
	return err
}

func (s *SprintService) AddTask(ctx context.Context, sprintID, taskID string) error {
	err := s.repo.AddTask(ctx, sprintID, taskID)
	if err != nil {
		s.logger.Error("service: add task to sprint", "error", err, "sprint_id", sprintID, "task_id", taskID)
	}
	return err
}

func (s *SprintService) RemoveTask(ctx context.Context, sprintID, taskID string) error {
	err := s.repo.RemoveTask(ctx, sprintID, taskID)
	if err != nil {
		s.logger.Error("service: remove task from sprint", "error", err, "sprint_id", sprintID, "task_id", taskID)
	}
	return err
}

package service

import (
	"context"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
)

type SprintService struct {
	repo     repository.SprintRepository
	activity repository.ActivityRepository
}

func NewSprintService(repo repository.SprintRepository, activity repository.ActivityRepository) *SprintService {
	return &SprintService{
		repo:     repo,
		activity: activity,
	}
}

func (s *SprintService) FindByProject(ctx context.Context, projectID string, activeOnly bool) ([]model.Sprint, error) {
	return s.repo.FindByProject(ctx, projectID, activeOnly)
}

func (s *SprintService) FindByID(ctx context.Context, id string) (*model.Sprint, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *SprintService) Create(ctx context.Context, sprint *model.Sprint) error {
	return s.repo.Create(ctx, sprint)
}

func (s *SprintService) Update(ctx context.Context, sprint *model.Sprint) error {
	return s.repo.Update(ctx, sprint)
}

func (s *SprintService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *SprintService) AddTask(ctx context.Context, sprintID, taskID string) error {
	return s.repo.AddTask(ctx, sprintID, taskID)
}

func (s *SprintService) RemoveTask(ctx context.Context, sprintID, taskID string) error {
	return s.repo.RemoveTask(ctx, sprintID, taskID)
}

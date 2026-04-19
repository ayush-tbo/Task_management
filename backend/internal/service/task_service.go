package service

import (
	"context"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
)

type TaskService struct {
	repo     repository.TaskRepository
	activity repository.ActivityRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) FindByID(ctx context.Context, id string) (*model.Task, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *TaskService) FindByProject(ctx context.Context, projectID string, filters repository.TaskFilters, page, pageSize int) ([]model.Task, int, error) {
	return s.repo.FindByProject(ctx, projectID, filters, page, pageSize)
}

func (s *TaskService) FindByAssignee(ctx context.Context, userID string, filters repository.TaskFilters, page, pageSize int) ([]model.Task, int, error) {
	return s.repo.FindByAssignee(ctx, userID, filters, page, pageSize)
}

func (s *TaskService) Create(ctx context.Context, task *model.Task) error {
	return s.repo.Create(ctx, task)
}

func (s *TaskService) Update(ctx context.Context, task *model.Task) error {
	return s.repo.Update(ctx, task)
}

func (s *TaskService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *TaskService) CountByStatus(ctx context.Context, projectID string) ([]model.StatusChartEntry, error) {
	return s.repo.CountByStatus(ctx, projectID)
}

func (s *TaskService) CountByPriority(ctx context.Context, projectID string) ([]model.PriorityChartEntry, error) {
	return s.repo.CountByPriority(ctx, projectID)
}

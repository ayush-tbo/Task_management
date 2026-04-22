package service

import (
	"context"
	"log/slog"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
)

type TaskService struct {
	repo     repository.TaskRepository
	activity repository.ActivityRepository
	logger   *slog.Logger
}

func NewTaskService(repo repository.TaskRepository, logger *slog.Logger) *TaskService {
	return &TaskService{
		repo:   repo,
		logger: logger,
	}
}

func (s *TaskService) FindByID(ctx context.Context, id string) (*model.Task, error) {
	task, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("service: find task by id", "error", err, "task_id", id)
	}
	return task, err
}

func (s *TaskService) FindByProject(ctx context.Context, projectID string, filters repository.TaskFilters, page, pageSize int) ([]model.Task, int, error) {
	tasks, total, err := s.repo.FindByProject(ctx, projectID, filters, page, pageSize)
	if err != nil {
		s.logger.Error("service: find tasks by project", "error", err, "project_id", projectID)
	}
	return tasks, total, err
}

func (s *TaskService) FindByAssignee(ctx context.Context, userID string, filters repository.TaskFilters, page, pageSize int) ([]model.Task, int, error) {
	tasks, total, err := s.repo.FindByAssignee(ctx, userID, filters, page, pageSize)
	if err != nil {
		s.logger.Error("service: find tasks by assignee", "error", err, "user_id", userID)
	}
	return tasks, total, err
}

func (s *TaskService) Create(ctx context.Context, task *model.Task) error {
	err := s.repo.Create(ctx, task)
	if err != nil {
		s.logger.Error("service: create task", "error", err, "task_id", task.ID, "project_id", task.ProjectID)
	}
	return err
}

func (s *TaskService) Update(ctx context.Context, task *model.Task) error {
	err := s.repo.Update(ctx, task)
	if err != nil {
		s.logger.Error("service: update task", "error", err, "task_id", task.ID)
	}
	return err
}

func (s *TaskService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("service: delete task", "error", err, "task_id", id)
	}
	return err
}

func (s *TaskService) CountByStatus(ctx context.Context, projectID string) ([]model.StatusChartEntry, error) {
	entries, err := s.repo.CountByStatus(ctx, projectID)
	if err != nil {
		s.logger.Error("service: count tasks by status", "error", err, "project_id", projectID)
	}
	return entries, err
}

func (s *TaskService) CountByPriority(ctx context.Context, projectID string) ([]model.PriorityChartEntry, error) {
	entries, err := s.repo.CountByPriority(ctx, projectID)
	if err != nil {
		s.logger.Error("service: count tasks by priority", "error", err, "project_id", projectID)
	}
	return entries, err
}

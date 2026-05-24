package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SprintService struct {
	repo        repository.SprintRepository
	taskRepo    repository.TaskRepository
	activity    repository.ActivityRepository
	mongoClient *mongo.Client
	logger      *slog.Logger
}

func NewSprintService(
	repo repository.SprintRepository,
	taskRepo repository.TaskRepository,
	activity repository.ActivityRepository,
	mongoClient *mongo.Client,
	logger *slog.Logger,
) *SprintService {
	return &SprintService{
		repo:        repo,
		taskRepo:    taskRepo,
		activity:    activity,
		mongoClient: mongoClient,
		logger:      logger,
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

func (s *SprintService) AddTaskToSprint(ctx context.Context, sprintID string, taskID string) error {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		task.SprintID = &sprintID
		task.UpdatedAt = time.Now()
		if err := s.taskRepo.Update(ctx, task); err != nil {
			return nil, err
		}
		if err := s.repo.AddTask(ctx, sprintID, task.ID); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		s.logger.Error("service: add task to sprint transaction", "error", err, "sprint_id", sprintID, "task_id", task.ID)
	}
	return err
}

func (s *SprintService) RemoveTaskFromSprint(ctx context.Context, sprintID string, taskID string) error {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		task.SprintID = nil
		task.UpdatedAt = time.Now()
		if err := s.taskRepo.Update(ctx, task); err != nil {
			return nil, err
		}
		if err := s.repo.RemoveTask(ctx, sprintID, task.ID); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		s.logger.Error("service: remove task from sprint transaction", "error", err, "sprint_id", sprintID, "task_id", task.ID)
	}
	return err
}

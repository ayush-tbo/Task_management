package repository

import (
	"context"

	"github.com/floqast/task-management/backend/internal/model"
)

type SprintRepository interface {
	FindByProject(ctx context.Context, projectID string, activeOnly bool) ([]model.Sprint, error)
	FindByID(ctx context.Context, id string) (*model.Sprint, error)
	Create(ctx context.Context, sprint *model.Sprint) error
	Update(ctx context.Context, sprint *model.Sprint) error
	Delete(ctx context.Context, id string) error
	AddTask(ctx context.Context, sprintID, taskID string) error
	RemoveTask(ctx context.Context, sprintID, taskID string) error
}

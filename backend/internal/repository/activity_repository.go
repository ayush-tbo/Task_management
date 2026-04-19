package repository

import (
	"context"

	"github.com/floqast/task-management/backend/internal/model"
)

type ActivityRepository interface {
	FindByProject(ctx context.Context, projectID string, action *model.ActivityAction, page, pageSize int) ([]model.ActivityEntry, int, error)
	FindByTask(ctx context.Context, taskID string, page, pageSize int) ([]model.ActivityEntry, int, error)
	Create(ctx context.Context, entry *model.ActivityEntry) error
}

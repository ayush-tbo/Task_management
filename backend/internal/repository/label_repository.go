package repository

import (
	"context"

	"github.com/floqast/task-management/backend/internal/model"
)

type LabelRepository interface {
	FindByProject(ctx context.Context, projectID string) ([]model.Label, error)
	FindByID(ctx context.Context, id string) (*model.Label, error)
	Create(ctx context.Context, label *model.Label) error
	Update(ctx context.Context, label *model.Label) error
	Delete(ctx context.Context, id string) error
}

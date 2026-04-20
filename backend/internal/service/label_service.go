package service

import (
	"context"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
)

type LabelService struct {
	repo     repository.LabelRepository
	activity repository.ActivityRepository
}

func NewLabelService(repo repository.LabelRepository, activity repository.ActivityRepository) *LabelService {
	return &LabelService{
		repo:     repo,
		activity: activity,
	}
}

func (s *LabelService) FindByProject(ctx context.Context, projectID string) ([]model.Label, error) {
	return s.repo.FindByProject(ctx, projectID)
}

func (s *LabelService) FindByID(ctx context.Context, id string) (*model.Label, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *LabelService) Create(ctx context.Context, label *model.Label) error {
	return s.repo.Create(ctx, label)
}

func (s *LabelService) Update(ctx context.Context, label *model.Label) error {
	return s.repo.Update(ctx, label)
}

func (s *LabelService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

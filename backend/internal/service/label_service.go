package service

import (
	"context"
	"log/slog"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
)

type LabelService struct {
	repo     repository.LabelRepository
	activity repository.ActivityRepository
	logger   *slog.Logger
}

func NewLabelService(repo repository.LabelRepository, activity repository.ActivityRepository, logger *slog.Logger) *LabelService {
	return &LabelService{
		repo:     repo,
		activity: activity,
		logger:   logger,
	}
}

func (s *LabelService) FindByProject(ctx context.Context, projectID string) ([]model.Label, error) {
	labels, err := s.repo.FindByProject(ctx, projectID)
	if err != nil {
		s.logger.Error("service: find labels by project", "error", err, "project_id", projectID)
	}
	return labels, err
}

func (s *LabelService) FindByID(ctx context.Context, id string) (*model.Label, error) {
	label, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("service: find label by id", "error", err, "label_id", id)
	}
	return label, err
}

func (s *LabelService) Create(ctx context.Context, label *model.Label) error {
	err := s.repo.Create(ctx, label)
	if err != nil {
		s.logger.Error("service: create label", "error", err, "label_id", label.ID, "project_id", label.ProjectID)
	}
	return err
}

func (s *LabelService) Update(ctx context.Context, label *model.Label) error {
	err := s.repo.Update(ctx, label)
	if err != nil {
		s.logger.Error("service: update label", "error", err, "label_id", label.ID)
	}
	return err
}

func (s *LabelService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("service: delete label", "error", err, "label_id", id)
	}
	return err
}

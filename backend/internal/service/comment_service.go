package service

import (
	"context"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
)

type CommentService struct {
	repo     repository.CommentRepository
	activity repository.ActivityRepository
}

func NewCommentService(repo repository.CommentRepository, activity repository.ActivityRepository) *CommentService {
	return &CommentService{
		repo:     repo,
		activity: activity,
	}
}

func (s *CommentService) FindByTask(ctx context.Context, taskID string) ([]model.Comment, error) {
	return s.repo.FindByTask(ctx, taskID)
}

func (s *CommentService) FindByID(ctx context.Context, id string) (*model.Comment, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *CommentService) Create(ctx context.Context, comment *model.Comment) error {
	return s.repo.Create(ctx, comment)
}

func (s *CommentService) Update(ctx context.Context, comment *model.Comment) error {
	return s.repo.Update(ctx, comment)
}

func (s *CommentService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

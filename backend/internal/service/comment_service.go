package service

import (
	"context"
	"log/slog"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
)

type CommentService struct {
	repo     repository.CommentRepository
	activity repository.ActivityRepository
	logger   *slog.Logger
}

func NewCommentService(repo repository.CommentRepository, activity repository.ActivityRepository, logger *slog.Logger) *CommentService {
	return &CommentService{
		repo:     repo,
		activity: activity,
		logger:   logger,
	}
}

func (s *CommentService) FindByTask(ctx context.Context, taskID string) ([]model.Comment, error) {
	comments, err := s.repo.FindByTask(ctx, taskID)
	if err != nil {
		s.logger.Error("service: find comments by task", "error", err, "task_id", taskID)
	}
	return comments, err
}

func (s *CommentService) FindByID(ctx context.Context, id string) (*model.Comment, error) {
	comment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("service: find comment by id", "error", err, "comment_id", id)
	}
	return comment, err
}

func (s *CommentService) Create(ctx context.Context, comment *model.Comment) error {
	err := s.repo.Create(ctx, comment)
	if err != nil {
		s.logger.Error("service: create comment", "error", err, "task_id", comment.TaskID, "user_id", comment.UserID)
	}
	return err
}

func (s *CommentService) Update(ctx context.Context, comment *model.Comment) error {
	err := s.repo.Update(ctx, comment)
	if err != nil {
		s.logger.Error("service: update comment", "error", err, "comment_id", comment.ID)
	}
	return err
}

func (s *CommentService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("service: delete comment", "error", err, "comment_id", id)
	}
	return err
}

func (s *CommentService) DeleteAll(ctx context.Context, taskID string) error {
	err := s.repo.DeleteAll(ctx, taskID)
	if err != nil {
		s.logger.Error("service: delete all comments", "error", err, "task_id", taskID)
	}
	return err
}

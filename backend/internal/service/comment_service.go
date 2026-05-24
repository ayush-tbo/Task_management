package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/model/dto"
	"github.com/floqast/task-management/backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentService struct {
	repo      repository.CommentRepository
	taskRepo  repository.TaskRepository
	activity  repository.ActivityRepository
	notifRepo repository.NotificationRepository
	logger    *slog.Logger
}

func NewCommentService(
	repo repository.CommentRepository,
	taskRepo repository.TaskRepository,
	activity repository.ActivityRepository,
	notifRepo repository.NotificationRepository,
	logger *slog.Logger,
) *CommentService {
	return &CommentService{
		repo:      repo,
		taskRepo:  taskRepo,
		activity:  activity,
		notifRepo: notifRepo,
		logger:    logger,
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

func (s *CommentService) CreateComment(ctx context.Context, user *model.User, taskID string, req dto.CreateCommentRequest) (*model.Comment, error) {
	comment := &model.Comment{
		ID:        primitive.NewObjectID().Hex(),
		TaskID:    taskID,
		UserID:    user.ID,
		User:      user,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.repo.Create(ctx, comment)
	if err != nil {
		s.logger.Error("service: create comment", "error", err, "task_id", taskID, "user_id", user.ID)
		return nil, err
	}

	go func() {
		bgCtx := context.Background()
		task, err := s.taskRepo.FindByID(bgCtx, taskID)
		if err == nil && task != nil {
			refType := "task"
			notified := map[string]bool{user.ID: true}
			notify := func(targetUserID string) {
				if notified[targetUserID] {
					return
				}
				notified[targetUserID] = true
				n := &model.Notification{
					ID:            primitive.NewObjectID().Hex(),
					UserID:        targetUserID,
					Type:          model.NotifMention,
					Title:         "New Comment",
					Message:       user.Name + " commented on \"" + task.Title + "\"",
					ReferenceType: &refType,
					ReferenceID:   &task.ID,
					CreatedAt:     time.Now(),
				}
				if err := s.notifRepo.Create(bgCtx, n); err != nil {
					s.logger.Error("notification create failed", "error", err)
				}
			}
			if task.ReporterID != "" {
				notify(task.ReporterID)
			}
			if task.AssigneeID != nil && *task.AssigneeID != "" {
				notify(*task.AssigneeID)
			}
		}

		entry := &model.ActivityEntry{
			ID:        primitive.NewObjectID().Hex(),
			ProjectID: req.ProjectID,
			TaskID:    &taskID,
			UserID:    user.ID,
			User:      user,
			Action:    model.ActionCommentAdded,
			Details:   map[string]interface{}{"comment_id": comment.ID},
			CreatedAt: time.Now(),
		}
		if err := s.activity.Create(bgCtx, entry); err != nil {
			s.logger.Error("activity log failed", "error", err)
		}
	}()

	return comment, nil
}

func (s *CommentService) UpdateComment(ctx context.Context, user *model.User, comment *model.Comment, req dto.UpdateCommentRequest) (*model.Comment, error) {
	oldContent := comment.Content

	if req.Content != nil {
		comment.Content = *req.Content
	}
	comment.UpdatedAt = time.Now()

	err := s.repo.Update(ctx, comment)
	if err != nil {
		s.logger.Error("service: update comment", "error", err, "comment_id", comment.ID)
		return nil, err
	}

	go func() {
		entry := &model.ActivityEntry{
			ID:        primitive.NewObjectID().Hex(),
			ProjectID: req.ProjectID,
			TaskID:    &req.TaskID,
			UserID:    user.ID,
			User:      user,
			Action:    model.ActionCommentChanged,
			Details:   map[string]interface{}{"comment_id": comment.ID, "old_content": oldContent},
			CreatedAt: time.Now(),
		}
		if err := s.activity.Create(context.Background(), entry); err != nil {
			s.logger.Error("activity log failed", "error", err)
		}
	}()

	return comment, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, user *model.User, comment *model.Comment, req dto.DeleteCommentRequest) error {
	err := s.repo.Delete(ctx, comment.ID)
	if err != nil {
		s.logger.Error("service: delete comment", "error", err, "comment_id", comment.ID)
		return err
	}

	go func() {
		entry := &model.ActivityEntry{
			ID:        primitive.NewObjectID().Hex(),
			ProjectID: req.ProjectID,
			TaskID:    &req.TaskID,
			UserID:    user.ID,
			User:      user,
			Action:    model.ActionCommentDeleted,
			Details:   map[string]interface{}{"comment_id": comment.ID, "info": "Comment was removed by the author"},
			CreatedAt: time.Now(),
		}
		if err := s.activity.Create(context.Background(), entry); err != nil {
			s.logger.Error("activity log failed", "error", err)
		}
	}()

	return nil
}

func (s *CommentService) DeleteAll(ctx context.Context, taskID string) error {
	err := s.repo.DeleteAll(ctx, taskID)
	if err != nil {
		s.logger.Error("service: delete all comments", "error", err, "task_id", taskID)
	}
	return err
}

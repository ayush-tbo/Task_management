package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/model/dto"
	"github.com/floqast/task-management/backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TaskService struct {
	repo        repository.TaskRepository
	commentRepo repository.CommentRepository
	projectRepo repository.ProjectRepository
	activity    repository.ActivityRepository
	notifRepo   repository.NotificationRepository
	mongoClient *mongo.Client
	logger      *slog.Logger
}

func NewTaskService(
	repo repository.TaskRepository,
	commentRepo repository.CommentRepository,
	projectRepo repository.ProjectRepository,
	activityRepo repository.ActivityRepository,
	notifRepo repository.NotificationRepository,
	mongoClient *mongo.Client,
	logger *slog.Logger,
) *TaskService {
	return &TaskService{
		repo:        repo,
		commentRepo: commentRepo,
		projectRepo: projectRepo,
		activity:    activityRepo,
		notifRepo:   notifRepo,
		mongoClient: mongoClient,
		logger:      logger,
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

func (s *TaskService) CreateTask(ctx context.Context, user *model.User, projectID string, req dto.CreateTaskRequest) (*model.Task, error) {
	status := req.Status
	if status == "" {
		status = model.StatusTodo
	}
	priority := req.Priority
	if priority == "" {
		priority = model.PriorityP3
	}

	task := &model.Task{
		ID:          primitive.NewObjectID().Hex(),
		ProjectID:   projectID,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		Priority:    priority,
		DueDate:     req.DueDate,
		AssigneeID:  req.AssigneeID,
		ReporterID:  user.ID,
		LabelIDs:    req.LabelIDs,
		SprintID:    req.SprintID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repo.Create(ctx, task)
	if err != nil {
		s.logger.Error("service: create task", "error", err, "project_id", projectID, "user_id", user.ID)
		return nil, err
	}

	_ = s.projectRepo.IncrementTaskCount(ctx, projectID, 1)

	go func() {
		bgCtx := context.Background()
		if task.AssigneeID != nil && *task.AssigneeID != "" && *task.AssigneeID != user.ID {
			refType := "task"
			n := &model.Notification{
				ID:            primitive.NewObjectID().Hex(),
				UserID:        *task.AssigneeID,
				Type:          model.NotifAssignment,
				Title:         "Task Assigned",
				Message:       user.Name + " assigned you to \"" + task.Title + "\"",
				ReferenceType: &refType,
				ReferenceID:   &task.ID,
				CreatedAt:     time.Now(),
			}
			if err := s.notifRepo.Create(bgCtx, n); err != nil {
				s.logger.Error("notification create failed", "error", err)
			}
		}

		entry := &model.ActivityEntry{
			ID:        primitive.NewObjectID().Hex(),
			ProjectID: projectID,
			TaskID:    &task.ID,
			UserID:    user.ID,
			User:      user,
			Action:    model.ActionTaskCreated,
			Details:   map[string]interface{}{"task_title": task.Title},
			CreatedAt: time.Now(),
		}
		if err := s.activity.Create(bgCtx, entry); err != nil {
			s.logger.Error("activity log failed", "error", err, "task_id", task.ID)
		}
	}()

	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, user *model.User, taskID string, req dto.UpdateTaskRequest) (*model.Task, error) {
	task, err := s.repo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, nil
	}

	oldAssigneeID := ""
	if task.AssigneeID != nil {
		oldAssigneeID = *task.AssigneeID
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if req.AssigneeID != nil {
		task.AssigneeID = req.AssigneeID
	}
	if req.LabelIDs != nil {
		task.LabelIDs = req.LabelIDs
	}
	if req.SprintID != nil {
		task.SprintID = req.SprintID
	}
	task.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, task)
	if err != nil {
		s.logger.Error("service: update task", "error", err, "task_id", taskID, "user_id", user.ID)
		return nil, err
	}

	go func() {
		bgCtx := context.Background()
		newAssigneeID := ""
		if task.AssigneeID != nil {
			newAssigneeID = *task.AssigneeID
		}
		if newAssigneeID != oldAssigneeID && newAssigneeID != "" && newAssigneeID != user.ID {
			refType := "task"
			n := &model.Notification{
				ID:            primitive.NewObjectID().Hex(),
				UserID:        newAssigneeID,
				Type:          model.NotifAssignment,
				Title:         "Task Assigned",
				Message:       user.Name + " assigned you to \"" + task.Title + "\"",
				ReferenceType: &refType,
				ReferenceID:   &task.ID,
				CreatedAt:     time.Now(),
			}
			if err := s.notifRepo.Create(bgCtx, n); err != nil {
				s.logger.Error("notification create failed", "error", err)
			}
		}

		entry := &model.ActivityEntry{
			ID:        primitive.NewObjectID().Hex(),
			ProjectID: task.ProjectID,
			TaskID:    &task.ID,
			UserID:    user.ID,
			User:      user,
			Action:    model.ActionTaskUpdated,
			Details:   map[string]interface{}{"task_title": task.Title},
			CreatedAt: time.Now(),
		}
		if err := s.activity.Create(bgCtx, entry); err != nil {
			s.logger.Error("activity log failed", "error", err, "task_id", task.ID)
		}
	}()

	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, user *model.User, taskID string) error {
	task, err := s.repo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return nil
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		s.logger.Error("service: start session for delete task", "error", err)
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		if err := s.repo.Delete(ctx, taskID); err != nil {
			return nil, err
		}
		if err := s.commentRepo.DeleteAll(ctx, taskID); err != nil {
			return nil, err
		}
		if err := s.projectRepo.IncrementTaskCount(ctx, task.ProjectID, -1); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		s.logger.Error("service: delete task transaction", "error", err, "task_id", taskID)
		return err
	}

	go func() {
		entry := &model.ActivityEntry{
			ID:        primitive.NewObjectID().Hex(),
			ProjectID: task.ProjectID,
			TaskID:    &taskID,
			UserID:    user.ID,
			User:      user,
			Action:    model.ActionTaskDeleted,
			Details:   map[string]interface{}{"task_title": task.Title},
			CreatedAt: time.Now(),
		}
		if err := s.activity.Create(context.Background(), entry); err != nil {
			s.logger.Error("activity log failed", "error", err, "task_id", taskID)
		}
	}()

	return nil
}

func (s *TaskService) AssignTask(ctx context.Context, user *model.User, taskID string, req dto.AssignTaskRequest) (*model.Task, error) {
	task, err := s.repo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, nil
	}

	if req.AssigneeID == "" {
		task.AssigneeID = nil
	} else {
		task.AssigneeID = &req.AssigneeID
	}
	task.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, task)
	if err != nil {
		s.logger.Error("service: assign task", "error", err, "task_id", taskID, "user_id", user.ID)
		return nil, err
	}

	go func() {
		bgCtx := context.Background()
		if req.AssigneeID != "" && req.AssigneeID != user.ID {
			refType := "task"
			n := &model.Notification{
				ID:            primitive.NewObjectID().Hex(),
				UserID:        req.AssigneeID,
				Type:          model.NotifAssignment,
				Title:         "Task Assigned",
				Message:       user.Name + " assigned you to \"" + task.Title + "\"",
				ReferenceType: &refType,
				ReferenceID:   &task.ID,
				CreatedAt:     time.Now(),
			}
			if err := s.notifRepo.Create(bgCtx, n); err != nil {
				s.logger.Error("notification create failed", "error", err)
			}
		}

		entry := &model.ActivityEntry{
			ID:        primitive.NewObjectID().Hex(),
			ProjectID: task.ProjectID,
			TaskID:    &task.ID,
			UserID:    user.ID,
			User:      user,
			Action:    model.ActionTaskAssigned,
			Details:   map[string]interface{}{"task_title": task.Title, "assignee_id": req.AssigneeID},
			CreatedAt: time.Now(),
		}
		if err := s.activity.Create(bgCtx, entry); err != nil {
			s.logger.Error("activity log", "error", err)
		}
	}()

	return task, nil
}

func (s *TaskService) UpdateStatus(ctx context.Context, user *model.User, taskID string, req dto.UpdateStatusRequest) (*model.Task, error) {
	task, err := s.repo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, nil
	}

	task.Status = req.Status
	task.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, task)
	if err != nil {
		s.logger.Error("service: update task status", "error", err, "task_id", taskID, "user_id", user.ID)
		return nil, err
	}

	go func() {
		bgCtx := context.Background()
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
				Type:          model.NotifAlert,
				Title:         "Status Updated",
				Message:       user.Name + " changed \"" + task.Title + "\" to " + string(req.Status),
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

		entry := &model.ActivityEntry{
			ID:        primitive.NewObjectID().Hex(),
			ProjectID: task.ProjectID,
			TaskID:    &task.ID,
			UserID:    user.ID,
			User:      user,
			Action:    model.ActionStatusChanged,
			Details:   map[string]interface{}{"task_title": task.Title, "new_status": string(req.Status)},
			CreatedAt: time.Now(),
		}
		if err := s.activity.Create(bgCtx, entry); err != nil {
			s.logger.Error("activity log", "error", err)
		}
	}()

	return task, nil
}

func (s *TaskService) LogTime(ctx context.Context, user *model.User, taskID string, req dto.LogTimeRequest) (*model.Task, error) {
	task, err := s.repo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, nil
	}

	if task.TimeTracking == nil {
		task.TimeTracking = &model.TimeTracking{}
	}
	task.TimeTracking.LoggedHours += req.Hours
	task.TimeTracking.Entries = append(task.TimeTracking.Entries, model.TimeEntry{
		Hours:       req.Hours,
		Description: req.Description,
		UserID:      user.ID,
		UserName:    user.Name,
		CreatedAt:   time.Now(),
	})
	task.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, task)
	if err != nil {
		s.logger.Error("service: log time", "error", err, "task_id", taskID)
		return nil, err
	}
	return task, nil
}

func (s *TaskService) CountByStatus(ctx context.Context, projectID string) ([]dto.StatusChartEntry, error) {
	entries, err := s.repo.CountByStatus(ctx, projectID)
	if err != nil {
		s.logger.Error("service: count tasks by status", "error", err, "project_id", projectID)
	}
	return entries, err
}

func (s *TaskService) CountByPriority(ctx context.Context, projectID string) ([]dto.PriorityChartEntry, error) {
	entries, err := s.repo.CountByPriority(ctx, projectID)
	if err != nil {
		s.logger.Error("service: count tasks by priority", "error", err, "project_id", projectID)
	}
	return entries, err
}

package repository

import (
	"context"

	"github.com/floqast/task-management/backend/internal/model"
)

type ProjectRepository interface {
	FindByID(ctx context.Context, id string) (*model.Project, error)
	FindByUser(ctx context.Context, userID string, page, pageSize int) ([]model.Project, int, error)
	Create(ctx context.Context, project *model.Project) error
	Update(ctx context.Context, project *model.Project) error
	Delete(ctx context.Context, id string) error
	ListMembers(ctx context.Context, projectID string) ([]model.ProjectMember, error)
	AddMember(ctx context.Context, projectID string, member *model.ProjectMember) error
	RemoveMember(ctx context.Context, projectID, userID string) error
}

type TaskRepository interface {
	FindByID(ctx context.Context, id string) (*model.Task, error)
	FindByProject(ctx context.Context, projectID string, filters TaskFilters, page, pageSize int) ([]model.Task, int, error)
	FindByAssignee(ctx context.Context, userID string, filters TaskFilters, page, pageSize int) ([]model.Task, int, error)
	Create(ctx context.Context, task *model.Task) error
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id string) error
	CountByStatus(ctx context.Context, projectID string) ([]model.StatusChartEntry, error)
	CountByPriority(ctx context.Context, projectID string) ([]model.PriorityChartEntry, error)
}

type TaskFilters struct {
	Status     *model.TaskStatus
	Priority   *model.Priority
	AssigneeID *string
	ReporterID *string
	LabelIDs   []string
	SprintID   *string
	DueBefore  *string
	DueAfter   *string
	IsPastDue  *bool
	SortBy     string
	SortOrder  string
}

// type CommentRepository interface {
// 	FindByTask(ctx context.Context, taskID string, page, pageSize int) ([]model.Comment, int, error)
// 	FindByID(ctx context.Context, id string) (*model.Comment, error)
// 	Create(ctx context.Context, comment *model.Comment) error
// 	Update(ctx context.Context, comment *model.Comment) error
// 	Delete(ctx context.Context, id string) error
// }

type LabelRepository interface {
	FindByProject(ctx context.Context, projectID string) ([]model.Label, error)
	FindByID(ctx context.Context, id string) (*model.Label, error)
	Create(ctx context.Context, label *model.Label) error
	Update(ctx context.Context, label *model.Label) error
	Delete(ctx context.Context, id string) error
}

type ActivityRepository interface {
	FindByProject(ctx context.Context, projectID string, action *model.ActivityAction, page, pageSize int) ([]model.ActivityEntry, int, error)
	FindByTask(ctx context.Context, taskID string, page, pageSize int) ([]model.ActivityEntry, int, error)
	Create(ctx context.Context, entry *model.ActivityEntry) error
}

type SprintRepository interface {
	FindByProject(ctx context.Context, projectID string, activeOnly bool) ([]model.Sprint, error)
	FindByID(ctx context.Context, id string) (*model.Sprint, error)
	Create(ctx context.Context, sprint *model.Sprint) error
	Update(ctx context.Context, sprint *model.Sprint) error
	Delete(ctx context.Context, id string) error
	AddTask(ctx context.Context, sprintID, taskID string) error
	RemoveTask(ctx context.Context, sprintID, taskID string) error
}

type NotificationRepository interface {
	FindByUser(ctx context.Context, userID string, unreadOnly bool, notifType *model.NotificationType, page, pageSize int) ([]model.Notification, int, int, error)
	MarkRead(ctx context.Context, id string) error
	MarkAllRead(ctx context.Context, userID string) error
	Create(ctx context.Context, notification *model.Notification) error
}

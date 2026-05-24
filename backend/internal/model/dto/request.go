package dto

import (
	"time"

	"github.com/floqast/task-management/backend/internal/model"
)

type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=200"`
	Description string `json:"description,omitempty" validate:"max=2000"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=200"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=2000"`
}

type AddMemberRequest struct {
	UserID string         `json:"user_id" validate:"required"`
	Role   model.MemberRole `json:"role" validate:"omitempty,oneof=admin member"`
}

type CreateTaskRequest struct {
	Title       string           `json:"title" validate:"required,min=1,max=500"`
	Description string           `json:"description,omitempty" validate:"max=10000"`
	Status      model.TaskStatus `json:"status,omitempty" validate:"omitempty,oneof=todo in_progress staging_review done"`
	Priority    model.Priority   `json:"priority,omitempty" validate:"omitempty,oneof=p1 p2 p3 p4"`
	DueDate     *time.Time       `json:"due_date,omitempty"`
	AssigneeID  *string          `json:"assignee_id,omitempty"`
	LabelIDs    []string         `json:"label_ids,omitempty"`
	SprintID    *string          `json:"sprint_id,omitempty"`
}

type UpdateTaskRequest struct {
	Title       *string           `json:"title,omitempty" validate:"omitempty,min=1,max=500"`
	Description *string           `json:"description,omitempty" validate:"omitempty,max=10000"`
	Status      *model.TaskStatus `json:"status,omitempty" validate:"omitempty,oneof=todo in_progress staging_review done"`
	Priority    *model.Priority   `json:"priority,omitempty" validate:"omitempty,oneof=p1 p2 p3 p4"`
	DueDate     *time.Time        `json:"due_date,omitempty"`
	AssigneeID  *string           `json:"assignee_id,omitempty"`
	LabelIDs    []string          `json:"label_ids,omitempty"`
	SprintID    *string           `json:"sprint_id,omitempty"`
}

type AssignTaskRequest struct {
	AssigneeID string `json:"assignee_id"`
}

type UpdateStatusRequest struct {
	Status model.TaskStatus `json:"status" validate:"required,oneof=todo in_progress staging_review done"`
}

type LogTimeRequest struct {
	Hours       float64 `json:"hours" validate:"required,gt=0"`
	Description string  `json:"description,omitempty" validate:"max=500"`
}

type CreateCommentRequest struct {
	Content   string `json:"content" validate:"required,min=1,max=5000"`
	ProjectID string `json:"project_id"`
}

type UpdateCommentRequest struct {
	Content   *string `json:"content" validate:"required,min=1,max=5000"`
	ProjectID string  `json:"project_id"`
	TaskID    string  `json:"task_id"`
}

type DeleteCommentRequest struct {
	ProjectID string `json:"project_id"`
	TaskID    string `json:"task_id"`
}

type CreateLabelRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=50"`
	Color string `json:"color,omitempty" validate:"omitempty,hexcolor"`
}

type CreateSprintRequest struct {
	Name      string    `json:"name" validate:"required,min=1,max=200"`
	Label     string    `json:"label,omitempty" validate:"max=100"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
}

type UpdateSprintRequest struct {
	Name      *string    `json:"name,omitempty" validate:"omitempty,min=1,max=200"`
	Label     *string    `json:"label,omitempty" validate:"omitempty,max=100"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	IsActive  *bool      `json:"is_active,omitempty"`
}

type RegisterUserRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Password  string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Name      *string `json:"name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty" validate:"omitempty,url"`
	Password  *string `json:"password,omitempty"`
}

type InviteUserRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	ProjectID *string `json:"project_id,omitempty"`
}

type SprintTaskRequest struct {
	TaskID string `json:"task_id" validate:"required"`
}

type CreateNotificationRequest struct {
	Title         string  `json:"title"`
	Message       string  `json:"message"`
	ReferenceType *string `json:"reference_type,omitempty"`
	ReferenceID   *string `json:"reference_id,omitempty"`
}

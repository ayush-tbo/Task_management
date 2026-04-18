package service

import (
	"github.com/floqast/task-management/backend/internal/repository"
)

type Services struct {
	User    *UserService
	Project *ProjectService
	Task    *TaskService
	// Comment      *CommentService
	Label        *LabelService
	Activity     *ActivityService
	Sprint       *SprintService
	Notification *NotificationService
}

type ProjectService struct {
	repo     repository.ProjectRepository
	activity repository.ActivityRepository
}

type TaskService struct {
	repo     repository.TaskRepository
	activity repository.ActivityRepository
}

// type CommentService struct {
// 	repo     repository.CommentRepository
// 	activity repository.ActivityRepository
// }

type LabelService struct {
	repo     repository.LabelRepository
	activity repository.ActivityRepository
}

type ActivityService struct {
	repo repository.ActivityRepository
}

type SprintService struct {
	repo     repository.SprintRepository
	activity repository.ActivityRepository
}

type NotificationService struct {
	repo repository.NotificationRepository
}

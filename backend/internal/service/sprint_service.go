package service

import "github.com/floqast/task-management/backend/internal/repository"

type SprintService struct {
	repo     repository.SprintRepository
	activity repository.ActivityRepository
}

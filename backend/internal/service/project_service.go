package service

import (
	"context"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
)

type ProjectService struct {
	repo     repository.ProjectRepository
	activity repository.ActivityRepository
}

func NewProjectService(repo repository.ProjectRepository, activityRepo repository.ActivityRepository) *ProjectService {
	return &ProjectService{
		repo:     repo,
		activity: activityRepo,
	}
}
func (s *ProjectService) Create(ctx context.Context, project *model.Project) error {
	return s.repo.Create(ctx, project)
}
func (s *ProjectService) FindByID(ctx context.Context, id string) (*model.Project, error) {
	return s.repo.FindByID(ctx, id)
}
func (s *ProjectService) Update(ctx context.Context, project *model.Project) error {
	return s.repo.Update(ctx, project)
}
func (s *ProjectService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
func (s *ProjectService) FindByUser(ctx context.Context, userID string, page, pageSize int) ([]model.Project, int, error) {
	return s.repo.FindByUser(ctx, userID, page, pageSize)
}
func (s *ProjectService) ListMembers(ctx context.Context, projectID string) ([]model.ProjectMember, error) {
	return s.repo.ListMembers(ctx, projectID)
}
func (s *ProjectService) AddMember(ctx context.Context, projectID string, member *model.ProjectMember) error {
	return s.repo.AddMember(ctx, projectID, member)
}
func (s *ProjectService) RemoveMember(ctx context.Context, projectID, userID string) error {
	return s.repo.RemoveMember(ctx, projectID, userID)
}

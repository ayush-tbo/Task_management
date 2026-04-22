package service

import (
	"context"
	"log/slog"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
)

type ProjectService struct {
	repo     repository.ProjectRepository
	activity repository.ActivityRepository
	logger   *slog.Logger
}

func NewProjectService(repo repository.ProjectRepository, activityRepo repository.ActivityRepository, logger *slog.Logger) *ProjectService {
	return &ProjectService{
		repo:     repo,
		activity: activityRepo,
		logger:   logger,
	}
}

func (s *ProjectService) Create(ctx context.Context, project *model.Project) error {
	err := s.repo.Create(ctx, project)
	if err != nil {
		s.logger.Error("service: create project", "error", err, "project_id", project.ID)
	}
	return err
}

func (s *ProjectService) FindByID(ctx context.Context, id string) (*model.Project, error) {
	project, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("service: find project by id", "error", err, "project_id", id)
	}
	return project, err
}

func (s *ProjectService) Update(ctx context.Context, project *model.Project) error {
	err := s.repo.Update(ctx, project)
	if err != nil {
		s.logger.Error("service: update project", "error", err, "project_id", project.ID)
	}
	return err
}

func (s *ProjectService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("service: delete project", "error", err, "project_id", id)
	}
	return err
}

func (s *ProjectService) FindByUser(ctx context.Context, userID string, page, pageSize int) ([]model.Project, int, error) {
	projects, total, err := s.repo.FindByUser(ctx, userID, page, pageSize)
	if err != nil {
		s.logger.Error("service: find projects by user", "error", err, "user_id", userID)
	}
	return projects, total, err
}

func (s *ProjectService) ListMembers(ctx context.Context, projectID string) ([]model.ProjectMember, error) {
	members, err := s.repo.ListMembers(ctx, projectID)
	if err != nil {
		s.logger.Error("service: list project members", "error", err, "project_id", projectID)
	}
	return members, err
}

func (s *ProjectService) AddMember(ctx context.Context, projectID string, member *model.ProjectMember) error {
	err := s.repo.AddMember(ctx, projectID, member)
	if err != nil {
		s.logger.Error("service: add project member", "error", err, "project_id", projectID, "user_id", member.UserID)
	}
	return err
}

func (s *ProjectService) RemoveMember(ctx context.Context, projectID, userID string) error {
	err := s.repo.RemoveMember(ctx, projectID, userID)
	if err != nil {
		s.logger.Error("service: remove project member", "error", err, "project_id", projectID, "user_id", userID)
	}
	return err
}

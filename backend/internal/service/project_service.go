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

type ProjectService struct {
	repo        repository.ProjectRepository
	userRepo    repository.UserRepository
	activity    repository.ActivityRepository
	notifRepo   repository.NotificationRepository
	mongoClient *mongo.Client
	logger      *slog.Logger
}

func NewProjectService(
	repo repository.ProjectRepository,
	userRepo repository.UserRepository,
	activityRepo repository.ActivityRepository,
	notifRepo repository.NotificationRepository,
	mongoClient *mongo.Client,
	logger *slog.Logger,
) *ProjectService {
	return &ProjectService{
		repo:        repo,
		userRepo:    userRepo,
		activity:    activityRepo,
		notifRepo:   notifRepo,
		mongoClient: mongoClient,
		logger:      logger,
	}
}

func (s *ProjectService) FindByID(ctx context.Context, id string) (*model.Project, error) {
	project, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("service: find project by id", "error", err, "project_id", id)
	}
	return project, err
}

func (s *ProjectService) FindByUser(ctx context.Context, userID string, page, pageSize int) ([]model.Project, int, error) {
	projects, total, err := s.repo.FindByUser(ctx, userID, page, pageSize)
	if err != nil {
		s.logger.Error("service: find projects by user", "error", err, "user_id", userID)
	}
	return projects, total, err
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

func (s *ProjectService) ListMembers(ctx context.Context, projectID string) ([]model.ProjectMember, error) {
	members, err := s.repo.ListMembers(ctx, projectID)
	if err != nil {
		s.logger.Error("service: list project members", "error", err, "project_id", projectID)
	}
	return members, err
}

func (s *ProjectService) CreateProject(ctx context.Context, user *model.User, req dto.CreateProjectRequest) (*model.Project, error) {
	project := &model.Project{
		ID:          primitive.NewObjectID().Hex(),
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     user.ID,
		MemberCount: 1,
		TaskCount:   0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	owner := &model.ProjectMember{
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		Role:      model.RoleOwner,
		JoinedAt:  time.Now(),
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		s.logger.Error("service: start session for create project", "error", err)
		return nil, err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		if err := s.repo.Create(ctx, project); err != nil {
			return nil, err
		}
		if err := s.repo.AddMember(ctx, project.ID, owner); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		s.logger.Error("service: create project transaction", "error", err, "user_id", user.ID)
		return nil, err
	}

	go func() {
		entry := &model.ActivityEntry{
			ID:        primitive.NewObjectID().Hex(),
			ProjectID: project.ID,
			UserID:    user.ID,
			User:      user,
			Action:    model.ActionMemberAdded,
			Details:   map[string]interface{}{"project_name": project.Name},
			CreatedAt: time.Now(),
		}
		if err := s.activity.Create(context.Background(), entry); err != nil {
			s.logger.Error("activity log failed", "error", err, "project_id", project.ID)
		}
	}()

	return project, nil
}

func (s *ProjectService) AddMember(ctx context.Context, user *model.User, projectID string, req dto.AddMemberRequest) (*model.ProjectMember, error) {
	role := req.Role
	if role == "" {
		role = model.RoleMember
	}

	userInfo, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil || userInfo == nil {
		s.logger.Error("service: find user for member", "error", err, "user_id", req.UserID)
		return nil, err
	}

	member := &model.ProjectMember{
		UserID:    req.UserID,
		Name:      userInfo.Name,
		Email:     userInfo.Email,
		AvatarURL: userInfo.AvatarURL,
		Role:      role,
		JoinedAt:  time.Now(),
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		if err := s.repo.AddMember(ctx, projectID, member); err != nil {
			return nil, err
		}
		if err := s.repo.IncrementMemberCount(ctx, projectID, 1); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		s.logger.Error("service: add member transaction", "error", err, "project_id", projectID)
		return nil, err
	}

	go func() {
		bgCtx := context.Background()
		project, _ := s.repo.FindByID(bgCtx, projectID)
		if req.UserID != user.ID {
			refType := "project"
			projectName := ""
			if project != nil {
				projectName = project.Name
			}
			n := &model.Notification{
				ID:            primitive.NewObjectID().Hex(),
				UserID:        req.UserID,
				Type:          model.NotifAlert,
				Title:         "Added to Project",
				Message:       user.Name + " added you to project \"" + projectName + "\"",
				ReferenceType: &refType,
				ReferenceID:   &projectID,
				CreatedAt:     time.Now(),
			}
			if err := s.notifRepo.Create(bgCtx, n); err != nil {
				s.logger.Error("notification create failed", "error", err)
			}
		}

		entry := &model.ActivityEntry{
			ID:        primitive.NewObjectID().Hex(),
			ProjectID: projectID,
			UserID:    user.ID,
			User:      user,
			Action:    model.ActionMemberAdded,
			Details:   map[string]interface{}{"member_user_id": req.UserID, "role": string(role)},
			CreatedAt: time.Now(),
		}
		if err := s.activity.Create(bgCtx, entry); err != nil {
			s.logger.Error("activity log", "error", err)
		}
	}()

	return member, nil
}

func (s *ProjectService) RemoveMember(ctx context.Context, user *model.User, projectID, userID string) error {
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		if err := s.repo.RemoveMember(ctx, projectID, userID); err != nil {
			return nil, err
		}
		if err := s.repo.IncrementMemberCount(ctx, projectID, -1); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		s.logger.Error("service: remove member transaction", "error", err, "project_id", projectID)
		return err
	}

	go func() {
		entry := &model.ActivityEntry{
			ID:        primitive.NewObjectID().Hex(),
			ProjectID: projectID,
			UserID:    user.ID,
			User:      user,
			Action:    model.ActionMemberRemoved,
			Details:   map[string]interface{}{"removed_user_id": userID},
			CreatedAt: time.Now(),
		}
		if err := s.activity.Create(context.Background(), entry); err != nil {
			s.logger.Error("activity log", "error", err)
		}
	}()

	return nil
}

func (s *ProjectService) IncrementTaskCount(ctx context.Context, projectID string, delta int) error {
	return s.repo.IncrementTaskCount(ctx, projectID, delta)
}

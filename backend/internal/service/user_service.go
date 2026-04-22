package service

import (
	"context"
	"log/slog"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
)

type UserService struct {
	repo   repository.UserRepository
	logger *slog.Logger
}

func NewUserService(repo repository.UserRepository, logger *slog.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

func UserIsAnonymous(u *model.User) bool {
	return u == model.AnonymousUser
}

func (s *UserService) FindByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("service: find user by id", "error", err, "user_id", id)
	}
	return user, err
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		s.logger.Error("service: find user by email", "error", err, "email", email)
	}
	return user, err
}

func (s *UserService) FindByName(ctx context.Context, name string) (*model.User, error) {
	user, err := s.repo.FindByName(ctx, name)
	if err != nil {
		s.logger.Error("service: find user by name", "error", err, "name", name)
	}
	return user, err
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		s.logger.Error("service: get all users", "error", err)
	}
	return users, err
}

func (s *UserService) Create(ctx context.Context, user *model.User) error {
	err := s.repo.Create(ctx, user)
	if err != nil {
		s.logger.Error("service: create user", "error", err, "user_id", user.ID, "email", user.Email)
	}
	return err
}

func (s *UserService) Update(ctx context.Context, user *model.User) error {
	err := s.repo.Update(ctx, user)
	if err != nil {
		s.logger.Error("service: update user", "error", err, "user_id", user.ID)
	}
	return err
}

func (s *UserService) GetToken(ctx context.Context, scope string, plainTextPassword string) (*model.Token, error) {
	token, err := s.repo.GetToken(ctx, scope, plainTextPassword)
	if err != nil {
		s.logger.Error("service: get token", "error", err, "scope", scope)
	}
	return token, err
}

func (s *UserService) CreateToken(ctx context.Context, token *model.Token) error {
	err := s.repo.CreateToken(ctx, token)
	if err != nil {
		s.logger.Error("service: create token", "error", err, "user_id", token.UserID)
	}
	return err
}

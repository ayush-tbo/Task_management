package service

import (
	"context"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func UserIsAnonymous(u *model.User) bool {
	return u == model.AnonymousUser
}

func (s *UserService) FindByID(ctx context.Context, id string) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *UserService) FindByName(ctx context.Context, name string) (*model.User, error) {
	return s.repo.FindByName(ctx, name)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.GetAllUsers(ctx)
}

func (s *UserService) Create(ctx context.Context, user *model.User) error {
	return s.repo.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, user *model.User) error {
	return s.repo.Update(ctx, user)
}

func (s *UserService) GetToken(ctx context.Context, scope string, plainTextPassword string) (*model.Token, error) {
	return s.repo.GetToken(ctx, scope, plainTextPassword)
}

func (s *UserService) CreateToken(ctx context.Context, token *model.Token) error {
	return s.repo.CreateToken(ctx, token)
}

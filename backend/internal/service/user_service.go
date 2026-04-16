package service

import (
	"context"
	"errors"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
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

func SetUserPassword(p *model.Password, plaintextpassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextpassword), 12)
	if err != nil {
		return err
	}

	p.PlainText = &plaintextpassword
	p.Hash = hash
	return nil
}

func MatchUserPassword(p *model.Password, plaintextpassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextpassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err // internal server error
		}
	}

	return true, nil
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

func (s *UserService) Create(ctx context.Context, user *model.User) error {
	return s.repo.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, user *model.User) error {
	return s.repo.Update(ctx, user)
}

func (s *UserService) GetUserToken(ctx context.Context, scope string, plainTextPassword string) (*model.User, error) {
	return s.repo.GetUserToken(ctx, scope, plainTextPassword)
}

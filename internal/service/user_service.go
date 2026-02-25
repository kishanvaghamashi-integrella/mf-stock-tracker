package service

import (
	"context"

	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/model"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/repository"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/util"
)

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, user *model.User) error {
	newHashedPassword, err := util.HashPassword(user.PasswordHash)
	if err != nil {
		return err
	}
	user.PasswordHash = newHashedPassword

	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, userId int) error {
	err := s.repo.Delete(ctx, userId)
	return err
}

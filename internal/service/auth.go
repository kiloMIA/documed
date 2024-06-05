package service

import (
	"context"

	"github.com/kiloMIA/documed/internal/entity"
	"github.com/kiloMIA/documed/internal/repo"
	"go.uber.org/zap"
)

type AuthService struct {
	userRepo repo.User
	logger   *zap.Logger
}

func NewAuthService(
	userRepo repo.User,
	logger *zap.Logger,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *AuthService) Create(ctx context.Context, input entity.CreateUser) error {
	user := entity.User{
		Name:     input.Name,
		Password: input.Password,
		Email:    input.Email,
	}

	err := s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.Error("auth service level - error creating user -", zap.Error(err))
		return err
	}

	return nil
}

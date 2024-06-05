package service

import (
	"context"

	"github.com/kiloMIA/documed/internal/entity"
)

type Auth interface {
	Create(ctx context.Context, input entity.CreateUser) error
}

type Service struct {
	Auth Auth
}

func NewService(authService Auth) *Service {
	return &Service{
		Auth: authService,
	}
}

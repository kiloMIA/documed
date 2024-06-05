package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kiloMIA/documed/internal/entity"
	"github.com/kiloMIA/documed/internal/repo/postgre"
	"go.uber.org/zap"
)

type User interface {
	Create(ctx context.Context, user entity.User) error
	Get(ctx context.Context, id int64) (entity.User, error)
	Update(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, id int64) error
}

type Repository struct {
	User
}

func NewRepository(dbpool *pgxpool.Pool, logger *zap.Logger) *Repository {
	return &Repository{
		User: postgre.NewUserRepo(dbpool, logger),
	}
}

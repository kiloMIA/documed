package postgre

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func CreateDB(logger *zap.Logger) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error("Failed to create pgxpool connection", zap.Error(err))
		return nil
	}

	return dbpool
}

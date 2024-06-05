package postgre

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kiloMIA/documed/internal/entity"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	dbpool *pgxpool.Pool
	logger *zap.Logger
}

func NewUserRepo(dbpool *pgxpool.Pool, logger *zap.Logger) *UserRepo {
	return &UserRepo{
		dbpool: dbpool,
		logger: logger,
	}
}

func (ur *UserRepo) Create(ctx context.Context, user entity.User) error {
	tx, err := ur.dbpool.Begin(ctx)
	if err != nil {
		ur.logger.Error("user repository level - error creating user - ", zap.Error(err))
		return err
	}
	defer tx.Rollback(ctx)

	hashedPassword := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	query := "INSERT INTO users (name, password, email) VALUES ($1, $2, $3)"
	_, err = tx.Exec(ctx, query, user.Name, string(hashedPassword), user.Email)
	if err != nil {
		ur.logger.Error("user repository level - error executing query - ", zap.Error(err))
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		ur.logger.Error("user repository level - error commiting transaction - ", zap.Error(err))
		return err
	}
	ur.logger.Info("user repository level - user succesfully created")
	return nil
}

func (ur *UserRepo) Get(ctx context.Context, id int64) (entity.User, error) {
	var user entity.User
	query := "SELECT id, name, email FROM users WHERE id = $1"
	row := ur.dbpool.QueryRow(ctx, query, id)

	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			ur.logger.Info("user repository level - no user found with id - ", zap.Int64("id", id))
			return user, nil
		}
		ur.logger.Error("user repository level - error scanning row - ", zap.Error(err))
		return user, err
	}
	ur.logger.Info("user repository level - user succesfully found")
	return user, nil
}

func (ur *UserRepo) Update(ctx context.Context, user entity.User) error {
	var hashedPassword string
	var err error

	if user.Password != "" {
		hashedPasswordBytes, err := bcrypt.GenerateFromPassword(
			[]byte(user.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			ur.logger.Error("user repository level - error hashing password -", zap.Error(err))
			return err
		}
		hashedPassword = string(hashedPasswordBytes)
	}

	tx, err := ur.dbpool.Begin(ctx)
	if err != nil {
		ur.logger.Error("user repository level - error starting transaction -", zap.Error(err))
		return err
	}
	defer tx.Rollback(ctx)

	query := `UPDATE users SET `
	args := []interface{}{}
	argIndex := 1

	if user.Name != "" {
		query += `name = $` + string(argIndex) + `, `
		args = append(args, user.Name)
		argIndex++
	}
	if user.Email != "" {
		query += `email = $` + string(argIndex) + `, `
		args = append(args, user.Email)
		argIndex++
	}
	if hashedPassword != "" {
		query += `password = $` + string(argIndex) + `, `
		args = append(args, hashedPassword)
		argIndex++
	}
	query = query[:len(query)-2]
	query += ` WHERE id = $` + string(argIndex)
	args = append(args, user.ID)

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		ur.logger.Error("user repository level - error executing query -", zap.Error(err))
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		ur.logger.Error("user repository level - error committing transaction -", zap.Error(err))
		return err
	}
	ur.logger.Info("user repository level - user successfully updated -", zap.Int64("id", user.ID))
	return nil
}

func (ur *UserRepo) Delete(ctx context.Context, id int64) error {
	tx, err := ur.dbpool.Begin(ctx)
	if err != nil {
		ur.logger.Error("user repository level - error starting transaction -", zap.Error(err))
		return err
	}
	defer tx.Rollback(ctx)

	sql := `DELETE FROM users WHERE id=$1`
	_, err = tx.Exec(ctx, sql, id)
	if err != nil {
		ur.logger.Error("user repository level - error executing query -", zap.Error(err))
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		ur.logger.Error("user repository level - error commiting transaction -", zap.Error(err))
		return err
	}
	ur.logger.Info("user repository level - user succesfully deleted -", zap.Int64("id", id))
	return nil
}

package repositoryimpl

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/model"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/util"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO USERS(name, email, password_hash)
		VALUES($1, $2, $3)
		RETURNING id
	`

	err := r.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, userId int64) error {
	query := `
		UPDATE users
		SET is_active=FALSE, updated_at=$2
		WHERE id=$1
	`

	res, err := r.db.Exec(ctx, query, userId, time.Now())
	if err != nil {
		slog.Error(err.Error())
		return util.NewInternalError("failed to delete user")
	}

	if res.RowsAffected() == 0 {
		return util.NewNotFoundError(fmt.Sprintf("no user found with id %d", userId))
	}
	return nil
}

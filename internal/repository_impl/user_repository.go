package repositoryimpl

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/model"
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

func (r *UserRepository) Get(ctx context.Context)    {}
func (r *UserRepository) Update(ctx context.Context) {}
func (r *UserRepository) Delete(ctx context.Context) {}

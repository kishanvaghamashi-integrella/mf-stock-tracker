package repository

import (
	"context"

	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/model"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, userId int64) error
}

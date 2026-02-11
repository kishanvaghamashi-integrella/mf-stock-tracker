package repository

import (
	"context"

	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/model"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context)
	Update(ctx context.Context)
	Delete(ctx context.Context)
}

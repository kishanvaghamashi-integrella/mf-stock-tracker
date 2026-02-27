package repository

import (
	"context"

	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/model"
)

type AssetRepositoryInterface interface {
	Create(ctx context.Context, asset *model.Asset) error
	GetByID(ctx context.Context, id int64) (*model.Asset, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.Asset, error)
	Update(ctx context.Context, asset *model.Asset) error
	Delete(ctx context.Context, id int64) error
}

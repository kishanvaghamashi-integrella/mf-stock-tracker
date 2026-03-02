package repository

import (
	"context"

	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/model"
)

type UserAssetRepositoryInterface interface {
	Create(ctx context.Context, userAsset *model.UserAsset) error
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.UserAsset, error)
	Delete(ctx context.Context, id int64) error
	IsUserAssetExits(ctx context.Context, userID int64, assetID int64) (bool, error)
}

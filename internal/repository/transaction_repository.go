package repository

import (
	"context"

	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/dto"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/model"
)

type TransactionRepositoryInterface interface {
	Create(ctx context.Context, txn *model.Transaction, holding *model.Holding, isUpdate bool) error
	GetAllByUserID(ctx context.Context, userID int64, limit, offset int) ([]dto.ResponseTransactionDto, error)
	GetHoldingsByUserAssetID(ctx context.Context, userAssetID int64) (*model.Holding, error)
	GetByID(ctx context.Context, id int64) (*model.Transaction, error)
	Update(ctx context.Context, txn *model.Transaction) error
	Delete(ctx context.Context, id int64) error
}

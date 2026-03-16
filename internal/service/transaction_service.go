package service

import (
	"context"
	"fmt"

	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/dto"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/model"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/repository"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/util"
)

type TransactionService struct {
	repo          repository.TransactionRepositoryInterface
	userAssetRepo repository.UserAssetRepositoryInterface
	userRepo      repository.UserRepositoryInterface
	assetRepo     repository.AssetRepositoryInterface
}

func NewTransactionService(
	repo repository.TransactionRepositoryInterface,
	userAssetRepo repository.UserAssetRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
	assetRepo repository.AssetRepositoryInterface,
) *TransactionService {
	return &TransactionService{repo: repo, userAssetRepo: userAssetRepo, userRepo: userRepo, assetRepo: assetRepo}
}

func (s *TransactionService) Create(ctx context.Context, req *dto.CreateTransactionRequest, userId int64) (*model.Transaction, error) {
	userExists, err := s.userRepo.ExistsByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, util.NewNotFoundError("User does not found on database")
	}

	assetExists, err := s.assetRepo.ExistsByID(ctx, req.AssetID)
	if err != nil {
		return nil, err
	}
	if !assetExists {
		return nil, util.NewNotFoundError("Asset does not found on database")
	}

	userAssetID, err := s.userAssetRepo.GetIdByUserIdAssetId(ctx, userId, req.AssetID)
	if err != nil {
		return nil, err
	}
	if userAssetID == nil {
		userAssetModel := &model.UserAsset{
			UserID:  userId,
			AssetID: req.AssetID,
		}
		if err := s.userAssetRepo.Create(ctx, userAssetModel); err != nil {
			return nil, err
		}
		userAssetID = &userAssetModel.ID
	}

	txn := &model.Transaction{
		UserAssetID: *userAssetID,
		TxnType:     req.TxnType,
		Quantity:    req.Quantity,
		Price:       req.Price,
		TxnDate:     req.TxnDate,
	}

	var isHoldingExisted bool = true
	holding, err := s.repo.GetHoldingsByUserAssetID(ctx, *userAssetID)
	if err != nil {
		return nil, err
	}
	if holding == nil {
		holding = &model.Holding{
			UserAssetID:   *userAssetID,
			TotalQuantity: txn.Quantity,
			AveragePrice:  txn.Price,
			TotalInvested: float64(txn.Price * txn.Quantity),
		}
		isHoldingExisted = false
	} else {
		if txn.TxnType == "BUY" {
			oldBoughtPrice := holding.AveragePrice
			oldBoughtQuantity := holding.TotalQuantity
			oldTotalPrice := oldBoughtPrice * oldBoughtQuantity

			newBoughtPrice := txn.Price
			newBoughtQuantity := txn.Quantity
			newTotalPrice := newBoughtPrice * newBoughtQuantity

			totalQuantity := oldBoughtQuantity + newBoughtQuantity
			totalPrice := oldTotalPrice + newTotalPrice
			totalAverage := totalPrice / totalQuantity

			holding.AveragePrice = totalAverage
			holding.TotalQuantity += newBoughtQuantity
			holding.TotalInvested = totalPrice
		} else {
			holding.TotalInvested -= txn.Price * txn.Quantity
			holding.TotalQuantity -= txn.Quantity
		}
	}

	if err := s.repo.Create(ctx, txn, holding, isHoldingExisted); err != nil {
		return nil, err
	}

	return txn, nil
}

func (s *TransactionService) GetAllByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.Transaction, error) {
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return nil, err
	}

	return s.repo.GetAllByUserID(ctx, userID, limit, offset)
}

func (s *TransactionService) Update(ctx context.Context, id int64, req *dto.UpdateTransactionRequest) error {
	txn, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if req.TxnType != nil {
		txn.TxnType = *req.TxnType
	}
	if req.Quantity != nil {
		txn.Quantity = *req.Quantity
	}
	if req.Price != nil {
		txn.Price = *req.Price
	}
	if req.TxnDate != nil {
		txn.TxnDate = *req.TxnDate
	}

	return s.repo.Update(ctx, txn)
}

func (s *TransactionService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *TransactionService) ensureUserExists(ctx context.Context, userID int64) error {
	exists, err := s.userRepo.ExistsByID(ctx, userID)
	if err != nil {
		return err
	}
	if !exists {
		return util.NewNotFoundError(fmt.Sprintf("user with id %d not found", userID))
	}
	return nil
}

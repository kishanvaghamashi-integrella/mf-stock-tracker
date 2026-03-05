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
}

func NewTransactionService(
	repo repository.TransactionRepositoryInterface,
	userAssetRepo repository.UserAssetRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
) *TransactionService {
	return &TransactionService{repo: repo, userAssetRepo: userAssetRepo, userRepo: userRepo}
}

func (s *TransactionService) Create(ctx context.Context, req *dto.CreateTransactionRequest) (*model.Transaction, error) {
	if err := s.ensureUserAssetExists(ctx, req.UserAssetID); err != nil {
		return nil, err
	}

	txn := &model.Transaction{
		UserAssetID: req.UserAssetID,
		TxnType:     req.TxnType,
		Quantity:    req.Quantity,
		Price:       req.Price,
		TxnDate:     req.TxnDate,
	}

	if err := s.repo.Create(ctx, txn); err != nil {
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

func (s *TransactionService) ensureUserAssetExists(ctx context.Context, userAssetID int64) error {
	exists, err := s.userAssetRepo.ExistsByID(ctx, userAssetID)
	if err != nil {
		return err
	}
	if !exists {
		return util.NewNotFoundError(fmt.Sprintf("user asset with id %d not found", userAssetID))
	}
	return nil
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

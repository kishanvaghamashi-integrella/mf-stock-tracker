package service

import (
	"context"
	"fmt"

	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/dto"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/repository"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/util"
)

type HoldingService struct {
	repo     repository.HoldingRepositoryInterface
	userRepo repository.UserRepositoryInterface
}

func NewHoldingService(
	repo repository.HoldingRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
) *HoldingService {
	return &HoldingService{repo: repo, userRepo: userRepo}
}

func (s *HoldingService) GetAllByUserID(ctx context.Context, userID int64, limit, offset int) ([]dto.HoldingResponseDto, error) {
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return nil, err
	}

	return s.repo.GetAllByUserID(ctx, userID, limit, offset)
}

func (s *HoldingService) ensureUserExists(ctx context.Context, userID int64) error {
	exists, err := s.userRepo.ExistsByID(ctx, userID)
	if err != nil {
		return err
	}
	if !exists {
		return util.NewNotFoundError(fmt.Sprintf("user with id %d not found", userID))
	}
	return nil
}

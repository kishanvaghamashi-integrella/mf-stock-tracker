package repository

import (
	"context"

	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/dto"
)

type HoldingRepositoryInterface interface {
	GetAllByUserID(ctx context.Context, userID int64, limit, offset int) ([]dto.HoldingResponseDto, error)
}

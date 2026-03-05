package dto

import "time"

type CreateTransactionRequest struct {
	UserAssetID int64     `json:"user_asset_id" validate:"required,gt=0"`
	TxnType     string    `json:"txn_type" validate:"required,txn_type"`
	Quantity    float64   `json:"quantity" validate:"required,gt=0"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	TxnDate     time.Time `json:"txn_date" validate:"required"`
}

type UpdateTransactionRequest struct {
	TxnType  *string  `json:"txn_type" validate:"omitempty,txn_type"`
	Quantity *float64 `json:"quantity" validate:"omitempty,gt=0"`
	Price    *float64 `json:"price" validate:"omitempty,gt=0"`
	TxnDate  *time.Time `json:"txn_date" validate:"omitempty"`
}

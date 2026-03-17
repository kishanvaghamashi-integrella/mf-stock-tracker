package dto

import "time"

type CreateTransactionRequest struct {
	AssetID  int64     `json:"asset_id" validate:"required,gt=0"`
	TxnType  string    `json:"txn_type" validate:"required,txn_type"`
	Quantity float64   `json:"quantity" validate:"required,gt=0"`
	Price    float64   `json:"price" validate:"required,gt=0"`
	TxnDate  time.Time `json:"txn_date" validate:"required"`
}

type UpdateTransactionRequest struct {
	ID       int64      `json:"id" validate:"required"`
	TxnType  *string    `json:"txn_type" validate:"omitempty,txn_type"`
	Quantity *float64   `json:"quantity" validate:"omitempty,gt=0"`
	Price    *float64   `json:"price" validate:"omitempty,gt=0"`
	TxnDate  *time.Time `json:"txn_date" validate:"omitempty"`
}

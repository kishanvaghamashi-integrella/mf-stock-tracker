package repositoryimpl

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/model"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/util"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, txn *model.Transaction) error {
	query := `
		INSERT INTO transactions (user_asset_id, txn_type, quantity, price, txn_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query, txn.UserAssetID, txn.TxnType, txn.Quantity, txn.Price, txn.TxnDate).
		Scan(&txn.ID, &txn.CreatedAt)
	if err != nil {
		slog.Error("failed to create transaction", "error", err.Error())
		return util.NewInternalError("failed to create transaction")
	}

	return nil
}

func (r *TransactionRepository) GetAllByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.Transaction, error) {
	query := `
		SELECT t.id, t.user_asset_id, t.txn_type, t.quantity, t.price, t.txn_date, t.created_at
		FROM transactions t
		INNER JOIN user_assets ua ON t.user_asset_id = ua.id
		WHERE ua.user_id = $1
		ORDER BY t.txn_date DESC, t.id DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		slog.Error("failed to list transactions", "error", err.Error())
		return nil, util.NewInternalError("failed to list transactions")
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var txn model.Transaction
		if err := rows.Scan(&txn.ID, &txn.UserAssetID, &txn.TxnType, &txn.Quantity, &txn.Price, &txn.TxnDate, &txn.CreatedAt); err != nil {
			slog.Error("failed to scan transaction row", "error", err.Error())
			return nil, util.NewInternalError("failed to list transactions")
		}
		transactions = append(transactions, txn)
	}

	if err := rows.Err(); err != nil {
		slog.Error("failed to iterate transaction rows", "error", err.Error())
		return nil, util.NewInternalError("failed to list transactions")
	}

	return transactions, nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, id int64) (*model.Transaction, error) {
	query := `
		SELECT id, user_asset_id, txn_type, quantity, price, txn_date, created_at
		FROM transactions
		WHERE id = $1
	`

	var txn model.Transaction
	err := r.db.QueryRow(ctx, query, id).
		Scan(&txn.ID, &txn.UserAssetID, &txn.TxnType, &txn.Quantity, &txn.Price, &txn.TxnDate, &txn.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, util.NewNotFoundError(fmt.Sprintf("transaction with id %d not found", id))
		}
		slog.Error("failed to get transaction", "error", err.Error())
		return nil, util.NewInternalError("failed to get transaction")
	}

	return &txn, nil
}

func (r *TransactionRepository) Update(ctx context.Context, txn *model.Transaction) error {
	query := `
		UPDATE transactions
		SET txn_type = $1, quantity = $2, price = $3, txn_date = $4
		WHERE id = $5
	`

	res, err := r.db.Exec(ctx, query, txn.TxnType, txn.Quantity, txn.Price, txn.TxnDate, txn.ID)
	if err != nil {
		slog.Error("failed to update transaction", "error", err.Error())
		return util.NewInternalError("failed to update transaction")
	}

	if res.RowsAffected() == 0 {
		return util.NewNotFoundError(fmt.Sprintf("transaction with id %d not found", txn.ID))
	}

	return nil
}

func (r *TransactionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM transactions WHERE id = $1`

	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		slog.Error("failed to delete transaction", "error", err.Error())
		return util.NewInternalError("failed to delete transaction")
	}

	if res.RowsAffected() == 0 {
		return util.NewNotFoundError(fmt.Sprintf("transaction with id %d not found", id))
	}

	return nil
}

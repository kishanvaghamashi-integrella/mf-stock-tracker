package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/dto"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/service"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/util"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// Create godoc
// @Summary Create transaction
// @Description Create a new transaction for a user asset
// @Tags transactions
// @Accept json
// @Produce json
// @Param payload body dto.CreateTransactionRequest true "Create transaction payload"
// @Success 201 {object} model.Transaction
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/transactions [post]
// @Security BearerAuth
func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	slog.Info("request started", "handler", "TransactionHandler.Create", "method", r.Method, "path", r.URL.Path)

	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("failed to decode request body", "handler", "TransactionHandler.Create", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := util.Validate.Struct(req); err != nil {
		slog.Warn("validation failed", "handler", "TransactionHandler.Create", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, util.FormatValidationErrors(err))
		return
	}

	txn, err := h.service.Create(r.Context(), &req)
	if err != nil {
		handleError(w, err, "TransactionHandler.Create")
		return
	}

	slog.Info("transaction created", "handler", "TransactionHandler.Create", "transactionID", txn.ID)
	util.SendResponse(w, http.StatusCreated, map[string]any{
		"message":     fmt.Sprintf("transaction created with id %d", txn.ID),
		"transaction": txn,
	})
}

// GetAllByUserID godoc
// @Summary List transactions for a user
// @Description Get all transactions across all user assets for a given user with pagination
// @Tags transactions
// @Produce json
// @Param limit query int false "Number of records to return (default: 50, max: 200)"
// @Param offset query int false "Number of records to skip (default: 0)"
// @Success 200 {array} model.Transaction
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/transactions [get]
// @Security BearerAuth
func (h *TransactionHandler) GetAllByUserID(w http.ResponseWriter, r *http.Request) {
	slog.Info("request started", "handler", "TransactionHandler.GetAllByUserID", "method", r.Method, "path", r.URL.Path)

	userID, ok := util.GetUserIDFromContext(r.Context())
	if !ok {
		slog.Warn("failed to parse user ID from context", "handler", "TransactionHandler.GetAllByUserID")
		util.SendErrorResponse(w, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	limit, offset, err := parsePaginationParams(r)
	if err != nil {
		slog.Warn("invalid pagination params", "handler", "TransactionHandler.GetAllByUserID", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	transactions, err := h.service.GetAllByUserID(r.Context(), userID, limit, offset)
	if err != nil {
		handleError(w, err, "TransactionHandler.GetAllByUserID")
		return
	}

	slog.Info("transactions retrieved", "handler", "TransactionHandler.GetAllByUserID", "userID", userID, "count", len(transactions), "limit", limit, "offset", offset)
	util.SendResponse(w, http.StatusOK, transactions)
}

// Update godoc
// @Summary Update transaction
// @Description Update an existing transaction by ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param txnId path int true "Transaction ID"
// @Param payload body dto.UpdateTransactionRequest true "Update transaction payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/transactions/{txnId} [put]
// @Security BearerAuth
func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	slog.Info("request started", "handler", "TransactionHandler.Update", "method", r.Method, "path", r.URL.Path)

	id, err := parseIntegerID(r, "txnId")
	if err != nil {
		slog.Warn("invalid transaction ID", "handler", "TransactionHandler.Update", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid transaction id")
		return
	}

	var req dto.UpdateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("failed to decode request body", "handler", "TransactionHandler.Update", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := util.Validate.Struct(req); err != nil {
		slog.Warn("validation failed", "handler", "TransactionHandler.Update", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, util.FormatValidationErrors(err))
		return
	}

	if err := h.service.Update(r.Context(), id, &req); err != nil {
		handleError(w, err, "TransactionHandler.Update")
		return
	}

	slog.Info("transaction updated", "handler", "TransactionHandler.Update", "transactionID", id)
	util.SendResponse(w, http.StatusOK, map[string]string{"message": "transaction updated successfully"})
}

// Delete godoc
// @Summary Delete transaction
// @Description Delete a transaction by its ID
// @Tags transactions
// @Produce json
// @Param txnId path int true "Transaction ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/transactions/{txnId} [delete]
// @Security BearerAuth
func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	slog.Info("request started", "handler", "TransactionHandler.Delete", "method", r.Method, "path", r.URL.Path)

	id, err := parseIntegerID(r, "txnId")
	if err != nil {
		slog.Warn("invalid transaction ID", "handler", "TransactionHandler.Delete", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid transaction id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		handleError(w, err, "TransactionHandler.Delete")
		return
	}

	slog.Info("transaction deleted", "handler", "TransactionHandler.Delete", "transactionID", id)
	util.SendResponse(w, http.StatusOK, map[string]string{"message": "transaction deleted successfully"})
}

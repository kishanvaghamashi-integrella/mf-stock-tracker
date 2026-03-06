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

type UserAssetHandler struct {
	service *service.UserAssetService
}

func NewUserAssetHandler(service *service.UserAssetService) *UserAssetHandler {
	return &UserAssetHandler{service: service}
}

// Create godoc
// @Summary Assign asset to user
// @Description Link an asset to a user
// @Tags user-assets
// @Accept json
// @Produce json
// @Param payload body dto.CreateUserAssetRequest true "Asset assignment payload"
// @Success 201 {object} model.UserAsset
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/user-assets [post]
// @Security BearerAuth
func (h *UserAssetHandler) Create(w http.ResponseWriter, r *http.Request) {
	slog.Info("request started", "handler", "UserAssetHandler.Create", "method", r.Method, "path", r.URL.Path)

	userID, ok := util.GetUserIDFromContext(r.Context())
	if ok == false {
		slog.Warn("failed to parse user ID from context", "handler", "UserAssetHandler.Create")
		util.SendErrorResponse(w, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	var req dto.CreateUserAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("failed to decode request body", "handler", "UserAssetHandler.Create", "userID", userID, "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := util.Validate.Struct(req); err != nil {
		slog.Warn("validation failed", "handler", "UserAssetHandler.Create", "userID", userID, "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, util.FormatValidationErrors(err))
		return
	}

	userAsset, err := h.service.Create(r.Context(), userID, &req)
	if err != nil {
		handleError(w, err, "UserAssetHandler.Create")
		return
	}

	slog.Info("user asset created", "handler", "UserAssetHandler.Create", "userID", userAsset.UserID, "assetID", userAsset.AssetID)
	util.SendResponse(w, http.StatusCreated, map[string]any{
		"message":    fmt.Sprintf("asset %d assigned to user %d", userAsset.AssetID, userAsset.UserID),
		"user_asset": userAsset,
	})
}

// GetByUserID godoc
// @Summary List assets for a user
// @Description Get all asset assignments for a user with pagination
// @Tags user-assets
// @Produce json
// @Param limit query int false "Number of records to return (default: 50, max: 200)"
// @Param offset query int false "Number of records to skip (default: 0)"
// @Success 200 {array} []model.UserAsset
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/user-assets [get]
// @Security BearerAuth
func (h *UserAssetHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	slog.Info("request started", "handler", "UserAssetHandler.GetByUserID", "method", r.Method, "path", r.URL.Path)

	userID, ok := util.GetUserIDFromContext(r.Context())
	if ok == false {
		slog.Warn("failed to parse user ID from context", "handler", "UserAssetHandler.GetByUserID")
		util.SendErrorResponse(w, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	limit, offset, err := parsePaginationParams(r)
	if err != nil {
		slog.Warn("invalid pagination params", "handler", "UserAssetHandler.GetByUserID", "userID", userID, "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userAssets, err := h.service.GetByUserID(r.Context(), userID, limit, offset)
	if err != nil {
		handleError(w, err, "UserAssetHandler.GetByUserID")
		return
	}

	slog.Info("user assets retrieved", "handler", "UserAssetHandler.GetByUserID", "userID", userID, "count", len(userAssets), "limit", limit, "offset", offset)
	util.SendResponse(w, http.StatusOK, userAssets)
}

// Delete godoc
// @Summary Remove asset assignment
// @Description Delete a user-asset link by its ID
// @Tags user-assets
// @Produce json
// @Param userAssetId path int true "User Asset ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/user-assets/{userAssetId} [delete]
// @Security BearerAuth
func (h *UserAssetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	slog.Info("request started", "handler", "UserAssetHandler.Delete", "method", r.Method, "path", r.URL.Path)

	userID, ok := util.GetUserIDFromContext(r.Context())
	if ok == false {
		slog.Warn("failed to parse user ID from context", "handler", "UserAssetHandler.Delete")
		util.SendErrorResponse(w, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	userAssetID, err := parseIntegerID(r, "userAssetId")
	if err != nil {
		slog.Warn("invalid user asset ID", "handler", "UserAssetHandler.Delete", "userID", userID, "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid user asset id")
		return
	}

	if err := h.service.Delete(r.Context(), userID, userAssetID); err != nil {
		handleError(w, err, "UserAssetHandler.Delete")
		return
	}

	slog.Info("user asset deleted", "handler", "UserAssetHandler.Delete", "userID", userID, "userAssetID", userAssetID)
	util.SendResponse(w, http.StatusOK, map[string]string{"message": "user asset deleted successfully"})
}

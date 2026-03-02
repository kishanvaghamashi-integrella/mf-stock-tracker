package handler

import (
	"encoding/json"
	"fmt"
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
// @Param userId path int true "User ID"
// @Param payload body dto.CreateUserAssetRequest true "Asset assignment payload"
// @Success 201 {object} model.UserAsset
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/users/{userId}/assets [post]
func (h *UserAssetHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := parseIntegerID(r, "userId")
	if err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req dto.CreateUserAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := util.Validate.Struct(req); err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, util.FormatValidationErrors(err))
		return
	}

	userAsset, err := h.service.Create(r.Context(), userID, &req)
	if err != nil {
		handleError(w, err)
		return
	}

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
// @Param userId path int true "User ID"
// @Param limit query int false "Number of records to return (default: 50, max: 200)"
// @Param offset query int false "Number of records to skip (default: 0)"
// @Success 200 {array} []model.UserAsset
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/users/{userId}/assets [get]
func (h *UserAssetHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := parseIntegerID(r, "userId")
	if err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid user id")
		return
	}

	limit, offset, err := parsePaginationParams(r)
	if err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userAssets, err := h.service.GetByUserID(r.Context(), userID, limit, offset)
	if err != nil {
		handleError(w, err)
		return
	}

	util.SendResponse(w, http.StatusOK, userAssets)
}

// Delete godoc
// @Summary Remove asset assignment
// @Description Delete a user-asset link by its ID
// @Tags user-assets
// @Produce json
// @Param userId path int true "User ID"
// @Param userAssetId path int true "User Asset ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/users/{userId}/assets/{userAssetId} [delete]
func (h *UserAssetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	_, err := parseIntegerID(r, "userId")
	if err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid user id")
		return
	}

	id, err := parseIntegerID(r, "userAssetId")
	if err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid user asset id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}

	util.SendResponse(w, http.StatusOK, map[string]string{"message": "user asset deleted successfully"})
}

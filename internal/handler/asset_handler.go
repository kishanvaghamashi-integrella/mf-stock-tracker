package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/dto"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/service"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/util"
)

type AssetHandler struct {
	service *service.AssetService
}

func NewAssetHandler(service *service.AssetService) *AssetHandler {
	return &AssetHandler{service: service}
}

// Create godoc
// @Summary Create asset
// @Description Create a new asset
// @Tags assets
// @Accept json
// @Produce json
// @Param payload body dto.CreateAssetRequest true "Create asset payload"
// @Success 201 {object} model.Asset
// @Failure 400 {object} util.ErrorBody
// @Failure 409 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/assets/ [post]
func (h *AssetHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createAssetDto dto.CreateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&createAssetDto); err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := util.Validate.Struct(createAssetDto); err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, util.FormatValidationErrors(err))
		return
	}

	asset, err := h.service.Create(r.Context(), &createAssetDto)
	if err != nil {
		handleError(w, err)
		return
	}

	util.SendResponse(w, http.StatusCreated, map[string]any{
		"message": fmt.Sprintf("asset created with id %d", asset.ID),
		"asset":   asset,
	})
}

// GetByID godoc
// @Summary Get asset by ID
// @Description Retrieve a single asset by its ID
// @Tags assets
// @Produce json
// @Param assetId path int64 true "Asset ID"
// @Success 200 {object} model.Asset
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/assets/{assetId} [get]
func (h *AssetHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIntegerID(r, "assetId")
	if err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid asset id")
		return
	}

	asset, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	util.SendResponse(w, http.StatusOK, asset)
}

// GetAll godoc
// @Summary List assets
// @Description Retrieve assets with pagination
// @Tags assets
// @Produce json
// @Param limit query int false "Number of records to return (default: 50, max: 200)"
// @Param offset query int false "Number of records to skip (default: 0)"
// @Success 200 {array} []model.Asset
// @Failure 400 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/assets/ [get]
func (h *AssetHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parsePaginationParams(r)
	if err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	assets, err := h.service.GetAll(r.Context(), limit, offset)
	if err != nil {
		handleError(w, err)
		return
	}

	util.SendResponse(w, http.StatusOK, assets)
}

// Update godoc
// @Summary Update asset
// @Description Update an existing asset by ID
// @Tags assets
// @Accept json
// @Produce json
// @Param assetId path int64 true "Asset ID"
// @Param payload body dto.UpdateAssetRequest true "Update asset payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/assets/{assetId} [put]
func (h *AssetHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIntegerID(r, "assetId")
	if err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid asset id")
		return
	}

	var updateAssetDto dto.UpdateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&updateAssetDto); err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := util.Validate.Struct(updateAssetDto); err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, util.FormatValidationErrors(err))
		return
	}

	if err := h.service.Update(r.Context(), id, &updateAssetDto); err != nil {
		handleError(w, err)
		return
	}

	util.SendResponse(w, http.StatusOK, map[string]string{"message": "asset updated successfully"})
}

// Delete godoc
// @Summary Delete asset
// @Description Delete an asset by ID
// @Tags assets
// @Produce json
// @Param assetId path int64 true "Asset ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/assets/{assetId} [delete]
func (h *AssetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIntegerID(r, "assetId")
	if err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid asset id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}

	util.SendResponse(w, http.StatusOK, map[string]string{"message": "asset deleted successfully"})
}

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
// @Security BearerAuth
func (h *AssetHandler) Create(w http.ResponseWriter, r *http.Request) {
	slog.Info("request started", "handler", "AssetHandler.Create", "method", r.Method, "path", r.URL.Path)

	var createAssetDto dto.CreateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&createAssetDto); err != nil {
		slog.Warn("failed to decode request body", "handler", "AssetHandler.Create", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := util.Validate.Struct(createAssetDto); err != nil {
		slog.Warn("validation failed", "handler", "AssetHandler.Create", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, util.FormatValidationErrors(err))
		return
	}

	asset, err := h.service.Create(r.Context(), &createAssetDto)
	if err != nil {
		util.HandleError(w, err, "AssetHandler.Create")
		return
	}

	slog.Info("asset created", "handler", "AssetHandler.Create", "assetID", asset.ID)
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
// @Security BearerAuth
func (h *AssetHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	slog.Info("request started", "handler", "AssetHandler.GetByID", "method", r.Method, "path", r.URL.Path)

	id, err := parseIntegerID(r, "assetId")
	if err != nil {
		slog.Warn("invalid asset ID", "handler", "AssetHandler.GetByID", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid asset id")
		return
	}

	asset, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		util.HandleError(w, err, "AssetHandler.GetByID")
		return
	}

	slog.Info("asset retrieved", "handler", "AssetHandler.GetByID", "assetID", id)
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
// @Security BearerAuth
func (h *AssetHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	slog.Info("request started", "handler", "AssetHandler.GetAll", "method", r.Method, "path", r.URL.Path)

	limit, offset, err := parsePaginationParams(r)
	if err != nil {
		slog.Warn("invalid pagination params", "handler", "AssetHandler.GetAll", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	assets, err := h.service.GetAll(r.Context(), limit, offset)
	if err != nil {
		util.HandleError(w, err, "AssetHandler.GetAll")
		return
	}

	slog.Info("assets retrieved", "handler", "AssetHandler.GetAll", "count", len(assets), "limit", limit, "offset", offset)
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
// @Security BearerAuth
func (h *AssetHandler) Update(w http.ResponseWriter, r *http.Request) {
	slog.Info("request started", "handler", "AssetHandler.Update", "method", r.Method, "path", r.URL.Path)

	id, err := parseIntegerID(r, "assetId")
	if err != nil {
		slog.Warn("invalid asset ID", "handler", "AssetHandler.Update", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid asset id")
		return
	}

	var updateAssetDto dto.UpdateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&updateAssetDto); err != nil {
		slog.Warn("failed to decode request body", "handler", "AssetHandler.Update", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := util.Validate.Struct(updateAssetDto); err != nil {
		slog.Warn("validation failed", "handler", "AssetHandler.Update", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, util.FormatValidationErrors(err))
		return
	}

	if err := h.service.Update(r.Context(), id, &updateAssetDto); err != nil {
		util.HandleError(w, err, "AssetHandler.Update")
		return
	}

	slog.Info("asset updated", "handler", "AssetHandler.Update", "assetID", id)
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
// @Security BearerAuth
func (h *AssetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	slog.Info("request started", "handler", "AssetHandler.Delete", "method", r.Method, "path", r.URL.Path)

	id, err := parseIntegerID(r, "assetId")
	if err != nil {
		slog.Warn("invalid asset ID", "handler", "AssetHandler.Delete", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid asset id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		util.HandleError(w, err, "AssetHandler.Delete")
		return
	}

	slog.Info("asset deleted", "handler", "AssetHandler.Delete", "assetID", id)
	util.SendResponse(w, http.StatusOK, map[string]string{"message": "asset deleted successfully"})
}

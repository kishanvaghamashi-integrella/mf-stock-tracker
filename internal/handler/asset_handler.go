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

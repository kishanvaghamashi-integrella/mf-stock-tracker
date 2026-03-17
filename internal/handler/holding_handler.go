package handler

import (
	"log/slog"
	"net/http"

	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/service"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/util"
)

type HoldingHandler struct {
	service *service.HoldingService
}

func NewHoldingHandler(service *service.HoldingService) *HoldingHandler {
	return &HoldingHandler{service: service}
}

// GetAll godoc
// @Summary List holdings for a user
// @Description Get all holdings for the authenticated user with pagination
// @Tags holdings
// @Produce json
// @Param limit query int false "Number of records to return (default: 50, max: 200)"
// @Param offset query int false "Number of records to skip (default: 0)"
// @Success 200 {array} dto.HoldingResponseDto
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/holdings [get]
// @Security BearerAuth
func (h *HoldingHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	slog.Info("request started", "handler", "HoldingHandler.GetAll", "method", r.Method, "path", r.URL.Path)

	userID, ok := util.GetUserIDFromContext(r.Context())
	if !ok {
		slog.Warn("failed to parse user ID from context", "handler", "HoldingHandler.GetAll")
		util.SendErrorResponse(w, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	limit, offset, err := parsePaginationParams(r)
	if err != nil {
		slog.Warn("invalid pagination params", "handler", "HoldingHandler.GetAll", "userID", userID, "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	holdings, err := h.service.GetAllByUserID(r.Context(), userID, limit, offset)
	if err != nil {
		util.HandleError(w, err, "HoldingHandler.GetAll")
		return
	}

	slog.Info("holdings retrieved", "handler", "HoldingHandler.GetAll", "userID", userID, "count", len(holdings), "limit", limit, "offset", offset)
	util.SendResponse(w, http.StatusOK, holdings)
}

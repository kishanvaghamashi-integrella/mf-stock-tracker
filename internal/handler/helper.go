package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/util"
)

func parseIntegerID(r *http.Request, param string) (int64, error) {
	return strconv.ParseInt(r.PathValue(param), 10, 64)
}

func handleError(w http.ResponseWriter, err error, handler string) {
	if appErr, ok := err.(*util.AppError); ok {
		if appErr.Code >= http.StatusInternalServerError {
			slog.Error("server error", "handler", handler, "status", appErr.Code, "error", appErr.Message)
		} else {
			slog.Warn("client error", "handler", handler, "status", appErr.Code, "error", appErr.Message)
		}
		util.SendErrorResponse(w, int(appErr.Code), appErr.Message)
	} else {
		slog.Error("unexpected error", "handler", handler, "error", err)
		util.SendErrorResponse(w, http.StatusInternalServerError, "unexpected error")
	}
}

func parsePaginationParams(r *http.Request) (int, int, error) {
	const (
		defaultLimit = 50
		maxLimit     = 200
	)

	limit := defaultLimit
	offset := 0

	if limitValue := r.URL.Query().Get("limit"); limitValue != "" {
		parsedLimit, err := strconv.Atoi(limitValue)
		if err != nil || parsedLimit <= 0 {
			return 0, 0, fmt.Errorf("limit must be a positive integer")
		}
		if parsedLimit > maxLimit {
			parsedLimit = maxLimit
		}
		limit = parsedLimit
	}

	if offsetValue := r.URL.Query().Get("offset"); offsetValue != "" {
		parsedOffset, err := strconv.Atoi(offsetValue)
		if err != nil || parsedOffset < 0 {
			return 0, 0, fmt.Errorf("offset must be a non-negative integer")
		}
		offset = parsedOffset
	}

	return limit, offset, nil
}

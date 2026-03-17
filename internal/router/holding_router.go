package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/handler"
)

func NewHoldingRouter(handler *handler.HoldingHandler) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handler.GetAll)

	return router
}

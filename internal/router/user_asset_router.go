package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/handler"
)

func NewUserAssetRouter(handler *handler.UserAssetHandler) http.Handler {
	router := chi.NewRouter()

	router.Post("/", handler.Create)
	router.Get("/", handler.GetByUserID)
	router.Delete("/{userAssetId}", handler.Delete)

	return router
}

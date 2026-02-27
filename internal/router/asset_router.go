package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/handler"
)

func NewAssetRouter(handler *handler.AssetHandler) http.Handler {
	router := chi.NewRouter()

	router.Post("/", handler.Create)
	router.Get("/", handler.GetAll)
	router.Get("/{assetId}", handler.GetByID)
	router.Put("/{assetId}", handler.Update)
	router.Delete("/{assetId}", handler.Delete)

	return router
}

package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/handler"
)

func NewUserRouter(handler *handler.UserHandler) http.Handler {
	router := chi.NewRouter()

	router.Post("/", handler.Create)
	router.Delete("/{userId}", handler.Delete)

	return router
}

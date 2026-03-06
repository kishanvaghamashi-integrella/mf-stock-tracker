package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/handler"
)

func NewTransactionRouter(handler *handler.TransactionHandler) http.Handler {
	router := chi.NewRouter()

	router.Post("/", handler.Create)
	router.Get("/", handler.GetAllByUserID)
	router.Put("/{txnId}", handler.Update)
	router.Delete("/{txnId}", handler.Delete)

	return router
}

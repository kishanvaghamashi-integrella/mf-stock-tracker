package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/handler"
	repositoryimpl "github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/repository_impl"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/router"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/service"
)

type App struct {
	Router http.Handler
}

func NewServer(db *pgxpool.Pool) *App {
	// Repository
	userRepo := repositoryimpl.NewUserRepository(db)
	assetRepo := repositoryimpl.NewAssetRepository(db)

	// Service
	userService := service.NewUserService(userRepo)
	assetService := service.NewAssetService(assetRepo)

	// Handler
	userHandler := handler.NewUserService(userService)
	assetHandler := handler.NewAssetHandler(assetService)

	r := chi.NewRouter()
	r.Mount("/api/users", router.NewUserRouter(userHandler))
	r.Mount("/api/assets", router.NewAssetRouter(assetHandler))
	return &App{Router: r}
}

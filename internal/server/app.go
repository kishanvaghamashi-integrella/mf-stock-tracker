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

	// Service
	userService := service.NewUserService(userRepo)

	// Handler
	userHandler := handler.NewUserService(userService)

	r := chi.NewRouter()
	r.Mount("/api/users", router.NewUserRouter(userHandler))
	return &App{Router: r}
}

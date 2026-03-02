package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/kishanvaghamashi-integrella/mf-stock-tracker/docs"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/handler"
	repositoryimpl "github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/repository_impl"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/router"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
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
	if isDevelopmentEnvironment() {
		r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/swagger/index.html", http.StatusTemporaryRedirect)
		})
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
	}
	r.Mount("/api/users", router.NewUserRouter(userHandler))
	r.Mount("/api/assets", router.NewAssetRouter(assetHandler))
	return &App{Router: r}
}

func isDevelopmentEnvironment() bool {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	return env == "dev" || env == "development" || env == "local"
}

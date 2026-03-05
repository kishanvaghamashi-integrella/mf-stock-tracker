package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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
	userAssetRepo := repositoryimpl.NewUserAssetRepository(db)
	txnRepo := repositoryimpl.NewTransactionRepository(db)

	// Service
	userService := service.NewUserService(userRepo)
	assetService := service.NewAssetService(assetRepo)
	userAssetService := service.NewUserAssetService(userAssetRepo, userRepo, assetRepo)
	txnService := service.NewTransactionService(txnRepo, userAssetRepo, userRepo)

	// Handler
	userHandler := handler.NewUserService(userService)
	assetHandler := handler.NewAssetHandler(assetService)
	userAssetHandler := handler.NewUserAssetHandler(userAssetService)
	txnHandler := handler.NewTransactionHandler(txnService)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	if isDevelopmentEnvironment() {
		r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/swagger/index.html", http.StatusTemporaryRedirect)
		})
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
	}
	r.Mount("/api/users", router.NewUserRouter(userHandler))
	r.Mount("/api/assets", router.NewAssetRouter(assetHandler))
	r.Mount("/api/users/{userId}/assets", router.NewUserAssetRouter(userAssetHandler))
	r.Mount("/api/transactions", router.NewTransactionRouter(txnHandler))
	return &App{Router: r}
}

func isDevelopmentEnvironment() bool {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	return env == "dev" || env == "development" || env == "local"
}

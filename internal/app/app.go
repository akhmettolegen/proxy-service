package app

import (
	"github.com/akhmettolegen/proxy-service/internal/config"
	v1 "github.com/akhmettolegen/proxy-service/internal/handler/http/v1"
	"github.com/akhmettolegen/proxy-service/internal/usecase"
	"github.com/akhmettolegen/proxy-service/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// Run creates objects via constructors.
func Run() {
	cfg := config.NewConfig()
	l := logger.New(cfg.LogLevel)
}

// NewRouter -.
// Swagger spec:
// @title       Proxy API
// @description Proxy service
// @version     1.0
// @BasePath    /v1
func setupRouter(l logger.Interface, t usecase.Task) chi.Router {
	router := chi.NewRouter()
	// Options
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Mount("/swagger", httpSwagger.WrapHandler)
	router.Get("/health", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

	// TODO: Prometheus metrics

	// Routes
	router.Mount("/v1/task", v1.NewTaskHandler(t, l).Routes())

	return router
}

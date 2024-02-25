package app

import (
	"fmt"
	"github.com/akhmettolegen/proxy-service/internal/config"
	"github.com/akhmettolegen/proxy-service/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Run creates objects via constructors.
func Run() {
	cfg := config.NewConfig()
	l := logger.New(cfg.LogLevel)

	// Use case
	translationUseCase := usecase.New(
		repo.New(pg),
		webapi.New(),
	)

	// HTTP Server
	router := setupRouter(l, translationUseCase)
	httpServer := httpserver.New(router, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

// NewRouter -.
// Swagger spec:
// @title       Translator API
// @description Translation service
// @version     1.0
// @BasePath    /v1
func setupRouter(l logger.Interface, t usecase.Translation) chi.Router {
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
	router.Mount("/v1/translation", v1.NewTranslationRoutes(t, l).Routes())

	return router
}

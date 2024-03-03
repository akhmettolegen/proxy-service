package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/akhmettolegen/proxy-service/internal/config"
	"github.com/akhmettolegen/proxy-service/internal/entity"
	v1 "github.com/akhmettolegen/proxy-service/internal/handler/http/v1"
	"github.com/akhmettolegen/proxy-service/internal/repo"
	service "github.com/akhmettolegen/proxy-service/internal/service"
	"github.com/akhmettolegen/proxy-service/internal/usecase"
	"github.com/akhmettolegen/proxy-service/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Run() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	l := logger.New(cfg.Log.Level)

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	mStorage := map[string]entity.Task{}
	mu := &sync.RWMutex{}
	taskR := repo.New(mStorage, mu)

	httpCli := &http.Client{}
	serv := service.NewClient(serverCtx, httpCli)

	taskUC := usecase.New(taskR, serv)

	router := setupRouter(l, taskUC)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Port,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, 10*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				l.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			l.Error(fmt.Errorf("failed to stop server %v", err))

			return
		}
		serverStopCtx()
	}()

	l.Info(fmt.Sprintf("listening and serving on port %s", cfg.HTTPServer.Port))
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		l.Error(fmt.Errorf("failed to start server %v", err))
	}

	<-serverCtx.Done()
	l.Info("server stopped")
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

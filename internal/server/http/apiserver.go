package http

import (
	"context"
	"github.com/akhmettolegen/test-service/internal/managers"
	v1 "github.com/akhmettolegen/test-service/internal/resources/http"
	proxyv1 "github.com/akhmettolegen/test-service/internal/resources/http/proxy/v1"
	"github.com/akhmettolegen/test-service/internal/server/configs"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"log"
	"net"
	"net/http"
	"time"
)

type APIServer struct {
	Address   string
	masterCtx context.Context

	proxyManager    managers.ProxyManager
	idleConnsClosed chan struct{}
	IsTesting       bool
}

func NewAPIServer(ctx context.Context, cfg *configs.Config, opts ...APIServerOption) *APIServer {
	srv := &APIServer{
		Address:         cfg.ListenAddr,
		masterCtx:       ctx,
		idleConnsClosed: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

// getAllowedOrigins возвращает список хостов для C.O.R.S.
func allowedOrigins(testing bool) []string {
	if testing {
		return []string{"*"}
	}

	return []string{}
}

func (srv *APIServer) Run() error {
	const (
		ReadTimeOut  = 30 * time.Second
		WriteTimeOut = 30 * time.Second
	)
	s := &http.Server{
		Addr:         srv.Address,
		Handler:      srv.setupRouter(),
		ReadTimeout:  ReadTimeOut,
		WriteTimeout: WriteTimeOut,
		BaseContext:  func(_ net.Listener) context.Context { return srv.masterCtx },
	}

	go srv.GracefulShutdown(s)
	log.Printf("[INFO] serving HTTP on \"%s\"", srv.Address)

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (srv *APIServer) setupRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.NoCache)   // no-cache
	r.Use(middleware.RequestID) // вставляет request ID в контекст каждого запроса
	r.Use(middleware.Logger)    // логирует начало и окончание каждого запроса с указанием времени обработки
	r.Use(middleware.Recoverer) // управляемо обрабатывает паники и выдает stack trace при их возникновении
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins(srv.IsTesting),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Mount("/version", v1.VersionResource{Version: "version"}.Routes())
	r.Mount("/proxy", proxyv1.ProxyResource{ProxyManager: srv.proxyManager}.Routes())

	return r
}

func (srv *APIServer) GracefulShutdown(httpSrv *http.Server) {
	<-srv.masterCtx.Done()

	if err := httpSrv.Shutdown(context.Background()); err != nil {
		log.Printf("[ERROR] HTTP server Shutdown: %v", err)
	}

	log.Println("[INFO] HTTP server has processed all idle connections")
	close(srv.idleConnsClosed)
}

func (srv *APIServer) Wait() {
	<-srv.idleConnsClosed
}

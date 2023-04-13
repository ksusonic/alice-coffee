package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ksusonic/alice-coffee/cloud/config"
	"go.uber.org/zap"
)

type Server struct {
	mux    *chi.Mux
	logger *zap.Logger
	cfg    *config.ServerConfig
}

type Route struct {
	Pattern string
	Handler http.Handler
}

func NewServer(cfg *config.ServerConfig, logger *zap.Logger, routers ...Route) *Server {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)

	for _, r := range routers {
		mux.Mount(r.Pattern, r.Handler)
	}

	return &Server{
		mux:    mux,
		logger: logger,
		cfg:    cfg,
	}
}

func (s Server) Run() *http.Server {
	s.logger.Info("Starting server", zap.String("address", s.cfg.Address))
	srv := &http.Server{
		Addr:    s.cfg.Address,
		Handler: s.mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Could not start srv listener", zap.Error(err))
		}
	}()
	return srv
}

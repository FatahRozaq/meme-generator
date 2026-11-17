// Package server menyediakan HTTP server dengan graceful shutdown.
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nyxsr/meme-portal/api/internal/config"
)

const (
	readHeaderTimeout = 10 * time.Second
	idleTimeout       = 120 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type Server struct {
	http *http.Server
}

func New(cfg config.Config, handler http.Handler) *Server {
	return &Server{
		http: &http.Server{
			Addr:              fmt.Sprintf(":%d", cfg.Port),
			Handler:           handler,
			ReadHeaderTimeout: readHeaderTimeout,
			IdleTimeout:       idleTimeout,
		},
	}
}

func (s *Server) Start() error {
	if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()
	return s.http.Shutdown(ctx)
}

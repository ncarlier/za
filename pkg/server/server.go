package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ncarlier/za/pkg/api"
	"github.com/ncarlier/za/pkg/config"
)

// Server instance
type Server struct {
	self *http.Server
}

// ListenAndServe starts server
func (s *Server) ListenAndServe() error {
	slog.Info("server is ready to handle requests", "addr", s.self.Addr)
	if err := s.self.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown server and managed service
func (s *Server) Shutdown(ctx context.Context) error {
	s.self.SetKeepAlivesEnabled(false)
	return s.self.Shutdown(ctx)
}

// NewServer creates new server instance
func NewServer(conf *config.Config) *Server {
	server := &Server{}
	server.self = &http.Server{
		Addr:    conf.HTTP.ListenAddr,
		Handler: api.NewRouter(conf),
	}
	return server
}

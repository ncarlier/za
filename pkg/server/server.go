package server

import (
	"context"
	"net/http"

	"github.com/ncarlier/za/pkg/api"
	"github.com/ncarlier/za/pkg/config"
	"github.com/ncarlier/za/pkg/logger"
)

// Server instance
type Server struct {
	self *http.Server
}

// ListenAndServe starts server
func (s *Server) ListenAndServe() error {
	logger.Info.Println("server is ready to handle requests at", s.self.Addr)
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
func NewServer(flags *config.Flags) *Server {
	conf := config.NewConfig()
	err := conf.LoadConfig(flags.ConfigFile)
	if err != nil {
		logger.Error.Fatal("unable to load configuration: ", err)
	}

	server := &Server{}
	server.self = &http.Server{
		Addr:    flags.ListenAddr,
		Handler: api.NewRouter(conf),
	}
	return server
}

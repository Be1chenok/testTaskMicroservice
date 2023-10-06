package server

import (
	"context"
	"net/http"
	"time"

	"github.com/Be1chenok/testTaskMicroservice/internal/config"
)

type Server struct {
	httpServer http.Server
}

func New(conf *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: http.Server{
			Addr:           conf.Server.Host + ":" + conf.Server.Port,
			MaxHeaderBytes: 1024 * 1024, // 1 MB
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			Handler:        handler,
		},
	}
}

func (s *Server) Start() error {
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Shuthdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

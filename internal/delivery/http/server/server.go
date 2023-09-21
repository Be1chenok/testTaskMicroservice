package server

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/Be1chenok/testTaskMicroservice/internal/config"
)

type Server struct {
	httpServer http.Server
}

func New(conf *config.ServerConfig, hadler http.Handler) *Server {
	return &Server{
		httpServer: http.Server{
			Addr:           conf.Host + ":" + conf.Port,
			MaxHeaderBytes: 1024 * 1024, // 1 MB
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			Handler:        hadler,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shuthdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello")
}

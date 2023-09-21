package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Be1chenok/testTaskMicroservice/internal/config"
	"github.com/Be1chenok/testTaskMicroservice/internal/delivery/http/server"
)

func Run() {
	srvConf, err := config.SrvInit()
	if err != nil {
		log.Fatalf("Failed to initialize server configuration: %v", err)
	}

	srv := server.New(srvConf, nil)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	<-quit

	log.Print("Shuthing Down")

	if err := srv.Shuthdown(context.Background()); err != nil {
		log.Fatalf("Failed to shut down the server: %v", err)
	}
}

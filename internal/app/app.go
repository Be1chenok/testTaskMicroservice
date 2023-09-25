package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Be1chenok/testTaskMicroservice/internal/config"
	"github.com/Be1chenok/testTaskMicroservice/internal/delivery/http/handler"
	"github.com/Be1chenok/testTaskMicroservice/internal/delivery/http/server"
	"github.com/Be1chenok/testTaskMicroservice/internal/repository"
	"github.com/Be1chenok/testTaskMicroservice/internal/repository/postgres"
	"github.com/Be1chenok/testTaskMicroservice/internal/service"
)

func Run() {
	ctx := context.Background()

	conf, err := config.Init()
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	db, err := postgres.New(conf)
	if err != nil {
		log.Fatal(err)
	}

	//cache, err := redis.New(ctx, conf)
	//if err != nil {
	//	log.Fatalf(err.Error())
	//}

	repository := repository.New(db)
	service := service.New(repository)
	handler := handler.New(service)
	srv := server.New(conf, handler.InitRoutes())

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	<-quit

	log.Print("Shuthing Down")

	if err := srv.Shuthdown(ctx); err != nil {
		log.Fatalf("Failed to shut down the server: %v", err)
	}
}

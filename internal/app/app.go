package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Be1chenok/testTaskMicroservice/internal/config"
	"github.com/Be1chenok/testTaskMicroservice/internal/delivery/http/handler"
	"github.com/Be1chenok/testTaskMicroservice/internal/delivery/http/server"
	"github.com/Be1chenok/testTaskMicroservice/internal/repository"
	"github.com/Be1chenok/testTaskMicroservice/internal/repository/postgres"
	rdb "github.com/Be1chenok/testTaskMicroservice/internal/repository/redis"
	"github.com/Be1chenok/testTaskMicroservice/internal/service"
	"github.com/Be1chenok/testTaskMicroservice/pkg/hash"
)

func Run() {
	ctx, cancle := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancle()

	conf, err := config.Init()
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	client, err := rdb.New(ctx, conf)
	if err != nil {
		log.Fatalf("Failed to connect Redis: %v", err)
	}

	db, err := postgres.New(conf)
	if err != nil {
		log.Fatalf("Failed to connect Postgres: %v", err)
	}

	hasher := hash.NewSHA256Hasher(conf.UserPassword.Salt)

	repository := repository.New(db, client)
	service := service.New(repository, hasher, conf.Tokens.SigningKey, conf.Tokens.AccessTTL, conf.Tokens.RefreshTTL)
	handler := handler.New(service)
	srv := server.New(conf, handler.InitRoutes())

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start the server: %v", err)
		}
	}()

	log.Print("Works fine")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	<-quit

	log.Print("Shuthing Down")

	if err := srv.Shuthdown(ctx); err != nil {
		log.Fatalf("Failed to shut down the server: %v", err)
	}

	if err = db.Close(); err != nil {
		log.Fatalf("Failed to close postgres: %v", err)
	}

	if err = client.Close(); err != nil {
		log.Fatalf("Failed to close redis: %v", err)
	}
}

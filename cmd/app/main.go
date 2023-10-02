package main

import (
	"github.com/Be1chenok/testTaskMicroservice/internal/app"
)

// @title Test Task Microservice
// @version pre-omega
// @descriptiom API Server for Authorization Application

// @127.0.0.1:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	app.Run()
}

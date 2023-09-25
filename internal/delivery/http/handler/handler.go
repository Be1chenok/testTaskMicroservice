package handler

import (
	"net/http"

	"github.com/Be1chenok/testTaskMicroservice/internal/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	router := mux.NewRouter()
	router.Use()
	return router
}

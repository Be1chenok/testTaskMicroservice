package handler

import (
	"net/http"

	"github.com/Be1chenok/testTaskMicroservice/internal/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	service *service.Service
	userId  int
}

func New(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/register", h.register).Methods("POST")
	router.HandleFunc("/login", h.login).Methods("POST")
	secure := router.PathPrefix("/auth").Subrouter()
	secure.Use(h.userIdentity)
	secure.HandleFunc("/home", h.homePage).Methods("GET")
	return router
}

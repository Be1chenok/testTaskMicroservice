package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Be1chenok/testTaskMicroservice/internal/service"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/Be1chenok/testTaskMicroservice/docs"
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

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/swagger.json", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "application/json")
		http.ServeFile(resp, req, "./docs/swagger.json")
	})

	router.HandleFunc("/register", h.register).Methods("POST")
	router.HandleFunc("/login", h.login).Methods("POST")
	secure := router.PathPrefix("/auth").Subrouter()
	secure.Use(h.userAccessIdentity)
	secure.HandleFunc("/home", h.homePage).Methods("GET")
	secure.HandleFunc("/logout", h.logOut).Methods("GET")
	secure.HandleFunc("/fullLogout", h.fullLogOut).Methods("GET")
	secure.HandleFunc("/refresh", h.refresh).Methods("GET")
	return router
}

func newErrorResponse(resp http.ResponseWriter, statusCode int, message string) {
	response := map[string]interface{}{
		"message": message,
	}
	resp.Header().Set(contentType, applicationJson)
	resp.WriteHeader(statusCode)
	json.NewEncoder(resp).Encode(response)
}

package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
)

const (
	contentType     = "Content-Type"
	applicationJson = "application/json"
)

type signInInput struct {
	Username string
	Password string
}

func (h *Handler) signUp(resp http.ResponseWriter, req *http.Request) {
	var input domain.User

	bytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(bytes, &input); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := h.service.CreateUser(input)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id": id,
	}
	resp.Header().Set(contentType, applicationJson)
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(response)
}

func (h *Handler) signIn(resp http.ResponseWriter, req *http.Request) {
	var input signInInput

	bytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(bytes, &input); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := h.service.GenerateToken(input.Username, input.Password)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"token": token,
	}
	resp.Header().Set(contentType, applicationJson)
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(response)
}

func (h *Handler) homePage(resp http.ResponseWriter, req *http.Request) {
	id := req.Context().Value(userCtx).(string)

	response := map[string]interface{}{
		"id": id,
	}
	resp.Header().Set(contentType, applicationJson)
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(response)
}

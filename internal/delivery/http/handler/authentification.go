package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
)

type signInInput struct {
	Username string
	Password string
}

func (h *Handler) signUp(resp http.ResponseWriter, req *http.Request) {
	var user domain.User

	bytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
	}

	if err := json.Unmarshal(bytes, &user); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}

	userId, err := h.service.CreateUser(user)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(string(fmt.Sprint(userId))))
}

func (h *Handler) signIn(resp http.ResponseWriter, req *http.Request) {
	var input signInInput

	bytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
	}

	if err := json.Unmarshal(bytes, &input); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}

	token, err := h.service.GenerateToken(input.Username, input.Password)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(token))

}

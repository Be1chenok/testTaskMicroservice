package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
	"github.com/Be1chenok/testTaskMicroservice/internal/model"
)

const (
	contentType     = "Content-Type"
	applicationJson = "application/json"
)

// @Sumary Register
// @Tags Authentefication
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body domain.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {string} string "message"
// @Failure 500 {string} string "message"
// @Failure default {string} string "message"
// @Router /register [post]
func (h *Handler) register(resp http.ResponseWriter, req *http.Request) {
	var input domain.User

	bytes, err := io.ReadAll(req.Body)
	if err != nil {
		newErrorResponse(resp, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.Unmarshal(bytes, &input); err != nil {
		newErrorResponse(resp, http.StatusInternalServerError, err.Error())
		return
	}

	userId, err := h.service.SignUp(model.SignUpInput{
		Email:    input.Email,
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		newErrorResponse(resp, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"userId": userId,
	}
	resp.Header().Set(contentType, applicationJson)
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(response)
}

// @Sumary Login
// @Tags Authentefication
// @Description create account
// @ID log-in
// @Accept json
// @Produce json
// @Param input body signInInput true "account info"
// @Success 200 {string} string "token"
// @Failure 400,404 {string} string "message"
// @Failure 500 {string} string "message"
// @Failure default {string} string "message"
// @Router /login [post]
func (h *Handler) login(resp http.ResponseWriter, req *http.Request) {
	var input model.SignInInput

	bytes, err := io.ReadAll(req.Body)
	if err != nil {
		newErrorResponse(resp, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.Unmarshal(bytes, &input); err != nil {
		newErrorResponse(resp, http.StatusInternalServerError, err.Error())
		return
	}

	tokens, err := h.service.SignIn(context.Background(), model.SignInInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		newErrorResponse(resp, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"accesToken":   tokens.AccesToken,
		"refreshToken": tokens.RefreshToken,
	}
	resp.Header().Set(contentType, applicationJson)
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(response)
}

// @Sumary LogOut
// @Security ApiKeyAuth
// @Tags Authentefication
// @Description logout
// @ID log-out
// @Accept json
// @Produce json
// @Success 200 {string} string "token"
// @Failure 400,404 {string} string "message"
// @Failure 500 {string} string "message"
// @Failure default {string} string "message"
// @Router /auth/logout [get]
func (h *Handler) logOut(resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get(authorizationHeader)
	headerParts := strings.Split(header, " ")

	accessToken := headerParts[1]
	if err := h.service.SignOut(context.Background(), accessToken); err != nil {
		newErrorResponse(resp, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"accessToken":  "",
		"refreshToken": "",
	}
	resp.Header().Set(contentType, applicationJson)
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(response)
}

// @Sumary FullLogOut
// @Security ApiKeyAuth
// @Tags Authentefication
// @Description fullLogout
// @ID full-logout
// @Accept json
// @Produce json
// @Success 200 {string} string "token"
// @Failure 400,404 {string} string "message"
// @Failure 500 {string} string "message"
// @Failure default {string} string "message"
// @Router /auth/fullLogout [get]
func (h *Handler) fullLogOut(resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get(authorizationHeader)
	headerParts := strings.Split(header, " ")
	accessToken := headerParts[1]

	if err := h.service.FullSignOut(context.Background(), accessToken); err != nil {
		newErrorResponse(resp, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"accessToken":  "",
		"refreshToken": "",
	}
	resp.Header().Set(contentType, applicationJson)
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(response)
}

// @Sumary Refresh
// @Security ApiKeyAuth
// @Tags Authentefication
// @Description refresh
// @ID refresh
// @Accept json
// @Produce json
// @Success 200 {string} string "token"
// @Failure 400,404 {string} string "message"
// @Failure 500 {string} string "message"
// @Failure default {string} string "message"
// @Router /auth/refresh [get]
func (h *Handler) refresh(resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get(authorizationHeader)
	headerParts := strings.Split(header, " ")
	refreshToken := headerParts[1]

	tokens, err := h.service.RefreshTokens(context.Background(), refreshToken)
	if err != nil {
		newErrorResponse(resp, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"accessToken":  tokens.AccesToken,
		"refreshToken": tokens.RefreshToken,
	}
	resp.Header().Set(contentType, applicationJson)
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(response)
}

// @Sumary HomePage
// @Security ApiKeyAuth
// @Tags Authentefication
// @Description homePage
// @ID home-page
// @Accept json
// @Produce json
// @Success 200 {string} string "home page for user"
// @Failure 400,404 {string} string "message"
// @Failure 500 {string} string "message"
// @Failure default {string} string "message"
// @Router /auth/home [get]
func (h *Handler) homePage(resp http.ResponseWriter, req *http.Request) {
	response := map[string]interface{}{
		"home page for user:": h.userId,
	}
	resp.Header().Set(contentType, applicationJson)
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(response)
}

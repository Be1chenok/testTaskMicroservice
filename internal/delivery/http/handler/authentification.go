package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

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
// @Success 200 {object} userResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
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
	response := userResponse{
		UserId: userId,
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
// @Param input body model.SignInInput true "account info"
// @Success 200 {object} tokensResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
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

	ctx, cancle := context.WithTimeout(req.Context(), 200*time.Millisecond)
	defer cancle()

	tokens, err := h.service.SignIn(ctx, model.SignInInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		newErrorResponse(resp, http.StatusInternalServerError, err.Error())
		return
	}

	response := tokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
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
// @Success 200 {integer} 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/logout [get]
func (h *Handler) logOut(resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get(authorizationHeader)
	headerParts := strings.Split(header, " ")

	accessToken := headerParts[1]

	ctx, cancle := context.WithTimeout(req.Context(), 200*time.Millisecond)
	defer cancle()

	if err := h.service.SignOut(ctx, accessToken); err != nil {
		newErrorResponse(resp, http.StatusInternalServerError, err.Error())
		return
	}

	resp.Header().Set(contentType, applicationJson)
	resp.WriteHeader(http.StatusOK)
}

// @Sumary FullLogOut
// @Security ApiKeyAuth
// @Tags Authentefication
// @Description fullLogout
// @ID full-logout
// @Accept json
// @Produce json
// @Success 200 {integer} 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/fullLogout [get]
func (h *Handler) fullLogOut(resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get(authorizationHeader)
	headerParts := strings.Split(header, " ")
	accessToken := headerParts[1]

	ctx, cancle := context.WithTimeout(req.Context(), 200*time.Millisecond)
	defer cancle()

	if err := h.service.FullSignOut(ctx, accessToken); err != nil {
		newErrorResponse(resp, http.StatusInternalServerError, err.Error())
		return
	}

	resp.Header().Set(contentType, applicationJson)
	resp.WriteHeader(http.StatusOK)
}

// @Sumary Refresh
// @Security ApiKeyAuth
// @Tags Authentefication
// @Description refresh
// @ID refresh
// @Accept json
// @Produce json
// @Success 200 {object} tokensResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/refresh [get]
func (h *Handler) refresh(resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get(authorizationHeader)
	headerParts := strings.Split(header, " ")
	refreshToken := headerParts[1]

	ctx, cancle := context.WithTimeout(req.Context(), 200*time.Millisecond)
	defer cancle()

	tokens, err := h.service.RefreshTokens(ctx, refreshToken)
	if err != nil {
		newErrorResponse(resp, http.StatusInternalServerError, err.Error())
		return
	}

	response := tokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
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
// @Success 200 {object} homeResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/home [get]
func (h *Handler) homePage(resp http.ResponseWriter, req *http.Request) {
	userId := req.Context().Value("userId").(int)
	response := homeResponse{
		UserId: userId,
	}
	resp.Header().Set(contentType, applicationJson)
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(response)
}

package handler

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

type tokensResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type userResponse struct {
	UserId int `json:"userId"`
}

type homeResponse struct {
	UserId int `json:"userId"`
}

func newErrorResponse(resp http.ResponseWriter, statusCode int, message string) {
	err := errorResponse{
		Message: message,
	}
	resp.Header().Set(contentType, applicationJson)
	resp.WriteHeader(statusCode)
	json.NewEncoder(resp).Encode(err)
}

package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const authorizationHeader = "Authorization"

func (h *Handler) userAccessIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		token, err := h.getTokenFromHeader(req)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusUnauthorized)
			newErrorResponse(resp, http.StatusUnauthorized, err.Error())
			return
		}

		userId, err := h.service.ParseToken(context.Background(), token)
		if err != nil {
			newErrorResponse(resp, http.StatusUnauthorized, err.Error())
			return
		}

		h.userId = userId

		next.ServeHTTP(resp, req)
	})
}

func (h *Handler) getTokenFromHeader(req *http.Request) (string, error) {
	header := req.Header.Get(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("access token is empty")
	}
	return headerParts[1], nil
}

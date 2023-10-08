package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

const authorizationHeader = "Authorization"

func (h *Handler) userAccessIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		token, err := h.getTokenFromHeader(req)
		if err != nil {
			newErrorResponse(resp, http.StatusUnauthorized, err.Error())
			return
		}

		userId, err := h.service.ParseToken(context.Background(), token)
		if err != nil {
			newErrorResponse(resp, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(req.Context(), "userId", userId)

		next.ServeHTTP(resp, req.WithContext(ctx))
	})
}

func (h *Handler) getTokenFromHeader(req *http.Request) (string, error) {
	header := req.Header.Get(authorizationHeader)
	if header == "" {
		return "", fmt.Errorf("Empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", fmt.Errorf("Invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", fmt.Errorf("token is empty")
	}
	return headerParts[1], nil
}

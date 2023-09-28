package handler

import (
	"net/http"
	"strings"
)

const authorizationHeader = "Authorization"

func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		header := req.Header.Get(authorizationHeader)
		if header == "" {
			http.Error(resp, "empty auth header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(resp, "invalid auth header", http.StatusUnauthorized)
			return
		}

		if len(headerParts[1]) == 0 {
			http.Error(resp, "token is empty", http.StatusUnauthorized)
			return
		}

		userId, err := h.service.Authentification.ParseToken(headerParts[1])
		if err != nil {
			http.Error(resp, err.Error(), http.StatusUnauthorized)
			return
		}

		h.userId = userId

		next.ServeHTTP(resp, req)
	})
}

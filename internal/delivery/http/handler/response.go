package handler

import (
	"encoding/json"
	"net/http"
)

func newErrorResponse(resp http.ResponseWriter, statusCode int, message string) {
	response := map[string]interface{}{
		"message": message,
	}
	resp.Header().Set(contentType, applicationJson)
	resp.WriteHeader(statusCode)
	json.NewEncoder(resp).Encode(response)
}

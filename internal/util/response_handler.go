package util

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Error any `json:"error"`
}

func SendResponse(w http.ResponseWriter, statusCode int, responseData any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(responseData)
}

func SendErrorResponse(w http.ResponseWriter, statusCode int, errData any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errBody := &ErrorBody{Error: errData}

	json.NewEncoder(w).Encode(errBody)
}

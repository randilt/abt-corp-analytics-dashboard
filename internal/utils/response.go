package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// WriteJSONResponse writes a JSON response
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

// WriteErrorResponse writes an error JSON response
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    statusCode,
	}

	WriteJSONResponse(w, statusCode, response)
}

// WriteSuccessResponse writes a success JSON response
func WriteSuccessResponse(w http.ResponseWriter, data interface{}) {
	response := SuccessResponse{
		Success: true,
		Data:    data,
	}

	WriteJSONResponse(w, http.StatusOK, response)
}

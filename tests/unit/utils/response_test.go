package utils_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"analytics-dashboard-api/internal/utils"
)

func TestWriteJSONResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		data       interface{}
		wantStatus int
	}{
		{
			name:       "success response",
			statusCode: http.StatusOK,
			data:       map[string]string{"message": "success"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "created response",
			statusCode: http.StatusCreated,
			data:       map[string]int{"id": 123},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "string response",
			statusCode: http.StatusOK,
			data:       "simple string",
			wantStatus: http.StatusOK,
		},
		{
			name:       "nil data",
			statusCode: http.StatusOK,
			data:       nil,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			utils.WriteJSONResponse(recorder, tt.statusCode, tt.data)

			// Check status code
			if recorder.Code != tt.wantStatus {
				t.Errorf("WriteJSONResponse() status = %d, want %d", recorder.Code, tt.wantStatus)
			}

			// Check content type
			contentType := recorder.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("WriteJSONResponse() content-type = %s, want application/json", contentType)
			}

			// Check if response body is valid JSON
			var result interface{}
			if err := json.NewDecoder(recorder.Body).Decode(&result); err != nil {
				t.Errorf("WriteJSONResponse() produced invalid JSON: %v", err)
			}
		})
	}
}

func TestWriteErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
	}{
		{
			name:       "bad request",
			statusCode: http.StatusBadRequest,
			message:    "Invalid input",
		},
		{
			name:       "internal server error",
			statusCode: http.StatusInternalServerError,
			message:    "Something went wrong",
		},
		{
			name:       "not found",
			statusCode: http.StatusNotFound,
			message:    "Resource not found",
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			message:    "Access denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			utils.WriteErrorResponse(recorder, tt.statusCode, tt.message)

			// Check status code
			if recorder.Code != tt.statusCode {
				t.Errorf("WriteErrorResponse() status = %d, want %d", recorder.Code, tt.statusCode)
			}

			// Check content type
			contentType := recorder.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("WriteErrorResponse() content-type = %s, want application/json", contentType)
			}

			// Parse response body
			var errorResp utils.ErrorResponse
			if err := json.NewDecoder(recorder.Body).Decode(&errorResp); err != nil {
				t.Fatalf("WriteErrorResponse() produced invalid JSON: %v", err)
			}

			// Check error response fields
			if errorResp.Message != tt.message {
				t.Errorf("WriteErrorResponse() message = %s, want %s", errorResp.Message, tt.message)
			}

			if errorResp.Code != tt.statusCode {
				t.Errorf("WriteErrorResponse() code = %d, want %d", errorResp.Code, tt.statusCode)
			}

			expectedError := http.StatusText(tt.statusCode)
			if errorResp.Error != expectedError {
				t.Errorf("WriteErrorResponse() error = %s, want %s", errorResp.Error, expectedError)
			}
		})
	}
}

func TestWriteSuccessResponse(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "map data",
			data: map[string]interface{}{
				"id":   123,
				"name": "test",
			},
		},
		{
			name: "string data",
			data: "success message",
		},
		{
			name: "array data",
			data: []string{"item1", "item2", "item3"},
		},
		{
			name: "nil data",
			data: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			utils.WriteSuccessResponse(recorder, tt.data)

			// Check status code
			if recorder.Code != http.StatusOK {
				t.Errorf("WriteSuccessResponse() status = %d, want %d", recorder.Code, http.StatusOK)
			}

			// Check content type
			contentType := recorder.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("WriteSuccessResponse() content-type = %s, want application/json", contentType)
			}

			// Parse response body
			var successResp utils.SuccessResponse
			if err := json.NewDecoder(recorder.Body).Decode(&successResp); err != nil {
				t.Fatalf("WriteSuccessResponse() produced invalid JSON: %v", err)
			}

			// Check success response fields
			if !successResp.Success {
				t.Error("WriteSuccessResponse() success should be true")
			}

			// Note: Deep comparison of Data field would require reflection or type assertion
			// For simplicity, we're just checking that the response structure is correct
		})
	}
}

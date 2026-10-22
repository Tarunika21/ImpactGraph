package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAnalysisHandlerRejectsInvalidRequests(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
	}{
		{
			name:           "Reject GET method",
			method:         http.MethodGet,
			body:           "",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Reject invalid JSON",
			method:         http.MethodPost,
			body:           `{"serviceName":`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Reject missing service name",
			method:         http.MethodPost,
			body:           `{"oldVersion":"1.0.0","newVersion":"2.0.0","oldSpec":{"paths":{}},"newSpec":{"paths":{}}}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Reject missing versions",
			method:         http.MethodPost,
			body:           `{"serviceName":"user-service","oldSpec":{"paths":{}},"newSpec":{"paths":{}}}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Reject missing specifications",
			method:         http.MethodPost,
			body:           `{"serviceName":"user-service","oldVersion":"1.0.0","newVersion":"2.0.0"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Reject unknown JSON fields",
			method:         http.MethodPost,
			body:           `{"serviceName":"user-service","unexpectedField":true}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Reject multiple JSON objects",
			method:         http.MethodPost,
			body:           `{} {}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Reject overly long service name",
			method:         http.MethodPost,
			body:           `{"serviceName":"` + strings.Repeat("a", 151) + `","oldVersion":"1.0.0","newVersion":"2.0.0","oldSpec":{"paths":{}},"newSpec":{"paths":{}}}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	handler := &AnalysisHandler{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(
				tt.method,
				"/api/v1/analyze",
				strings.NewReader(tt.body),
			)

			recorder := httptest.NewRecorder()

			handler.Handle(recorder, request)

			if recorder.Code != tt.expectedStatus {
				t.Errorf(
					"expected HTTP %d, got HTTP %d. Response: %s",
					tt.expectedStatus,
					recorder.Code,
					recorder.Body.String(),
				)
			}
		})
	}
}

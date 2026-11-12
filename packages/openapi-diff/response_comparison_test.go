package openapidiff

import (
	"testing"
)

func TestCompareResponseStatusCodes(t *testing.T) {
	tests := []struct {
		name         string
		oldResponses map[string]Response
		newResponses map[string]Response
		expectedType string
		expectedCode string
		breaking     bool
	}{
		{
			name: "removed 404 response",
			oldResponses: map[string]Response{
				"200": {},
				"404": {},
			},
			newResponses: map[string]Response{
				"200": {},
			},
			expectedType: "RESPONSE_STATUS_REMOVED",
			expectedCode: "404",
			breaking:     true,
		},
		{
			name: "added 201 response",
			oldResponses: map[string]Response{
				"200": {},
			},
			newResponses: map[string]Response{
				"200": {},
				"201": {},
			},
			expectedType: "RESPONSE_STATUS_ADDED",
			expectedCode: "201",
			breaking:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComparisonResponse{
				Changes: []Change{},
			}

			compareResponseStatusCodes(
				&result,
				"GET /users",
				tt.oldResponses,
				tt.newResponses,
			)

			if len(result.Changes) != 1 {
				t.Fatalf("expected 1 change, got %d", len(result.Changes))
			}

			change := result.Changes[0]

			if change.Type != tt.expectedType {
				t.Errorf("expected type %s, got %s", tt.expectedType, change.Type)
			}

			if change.Field != tt.expectedCode {
				t.Errorf("expected status %s, got %s", tt.expectedCode, change.Field)
			}

			if change.Breaking != tt.breaking {
				t.Errorf("expected breaking=%v, got %v", tt.breaking, change.Breaking)
			}
		})
	}
}

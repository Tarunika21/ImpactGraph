package openapidiff

import (
	"testing"
)

func TestCompareParameters(t *testing.T) {
	tests := []struct {
		name         string
		oldParams    []Parameter
		newParams    []Parameter
		expectedType string
		breaking     bool
	}{
		{
			name:      "new required parameter",
			oldParams: []Parameter{},
			newParams: []Parameter{
				{Name: "region", In: "query", Required: true, Schema: Schema{Type: "string"}},
			},
			expectedType: "PARAMETER_ADDED",
			breaking:     true,
		},
		{
			name:      "new optional parameter",
			oldParams: []Parameter{},
			newParams: []Parameter{
				{Name: "region", In: "query", Required: false, Schema: Schema{Type: "string"}},
			},
			expectedType: "PARAMETER_ADDED",
			breaking:     false,
		},
		{
			name: "parameter became required",
			oldParams: []Parameter{
				{Name: "region", In: "query", Required: false},
			},
			newParams: []Parameter{
				{Name: "region", In: "query", Required: true},
			},
			expectedType: "PARAMETER_BECAME_REQUIRED",
			breaking:     true,
		},
		{
			name: "parameter became optional",
			oldParams: []Parameter{
				{Name: "region", In: "query", Required: true},
			},
			newParams: []Parameter{
				{Name: "region", In: "query", Required: false},
			},
			expectedType: "PARAMETER_BECAME_OPTIONAL",
			breaking:     false,
		},
		{
			name: "parameter type changed",
			oldParams: []Parameter{
				{Name: "region", In: "query", Schema: Schema{Type: "string"}},
			},
			newParams: []Parameter{
				{Name: "region", In: "query", Schema: Schema{Type: "integer"}},
			},
			expectedType: "PARAMETER_TYPE_CHANGED",
			breaking:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComparisonResponse{
				Changes: []Change{},
			}

			compareParameters(&result, "GET /users", tt.oldParams, tt.newParams)

			if len(result.Changes) != 1 {
				t.Fatalf("expected 1 change, got %d", len(result.Changes))
			}

			change := result.Changes[0]

			if change.Type != tt.expectedType {
				t.Errorf("expected %s, got %s", tt.expectedType, change.Type)
			}

			if change.Breaking != tt.breaking {
				t.Errorf("expected breaking=%v, got %v", tt.breaking, change.Breaking)
			}
		})
	}
}

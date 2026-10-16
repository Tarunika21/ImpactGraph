package services

import (
	"testing"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
)

func TestCompareParameters(t *testing.T) {
	tests := []struct {
		name         string
		oldParams    []models.Parameter
		newParams    []models.Parameter
		expectedType string
		breaking     bool
	}{
		{
			name:      "new required parameter",
			oldParams: []models.Parameter{},
			newParams: []models.Parameter{
				{Name: "region", In: "query", Required: true, Schema: models.Schema{Type: "string"}},
			},
			expectedType: "PARAMETER_ADDED",
			breaking:     true,
		},
		{
			name:      "new optional parameter",
			oldParams: []models.Parameter{},
			newParams: []models.Parameter{
				{Name: "region", In: "query", Required: false, Schema: models.Schema{Type: "string"}},
			},
			expectedType: "PARAMETER_ADDED",
			breaking:     false,
		},
		{
			name: "parameter became required",
			oldParams: []models.Parameter{
				{Name: "region", In: "query", Required: false},
			},
			newParams: []models.Parameter{
				{Name: "region", In: "query", Required: true},
			},
			expectedType: "PARAMETER_BECAME_REQUIRED",
			breaking:     true,
		},
		{
			name: "parameter became optional",
			oldParams: []models.Parameter{
				{Name: "region", In: "query", Required: true},
			},
			newParams: []models.Parameter{
				{Name: "region", In: "query", Required: false},
			},
			expectedType: "PARAMETER_BECAME_OPTIONAL",
			breaking:     false,
		},
		{
			name: "parameter type changed",
			oldParams: []models.Parameter{
				{Name: "region", In: "query", Schema: models.Schema{Type: "string"}},
			},
			newParams: []models.Parameter{
				{Name: "region", In: "query", Schema: models.Schema{Type: "integer"}},
			},
			expectedType: "PARAMETER_TYPE_CHANGED",
			breaking:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := models.ComparisonResponse{
				Changes: []models.Change{},
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

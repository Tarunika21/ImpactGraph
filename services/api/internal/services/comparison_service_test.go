package services

import (
	"testing"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
)

func TestCompareSpecs(t *testing.T) {
	tests := []struct {
		name         string
		old          models.Schema
		new          models.Schema
		expectedType string
		expectedPath string
		breaking     bool
	}{
		{
			name: "field removed",
			old: models.Schema{Type: "object", Properties: map[string]models.Schema{
				"name": {Type: "string"},
			}},
			new:          models.Schema{Type: "object", Properties: map[string]models.Schema{}},
			expectedType: "FIELD_REMOVED",
			expectedPath: "name",
			breaking:     true,
		},
		{
			name: "field added",
			old:  models.Schema{Type: "object", Properties: map[string]models.Schema{}},
			new: models.Schema{Type: "object", Properties: map[string]models.Schema{
				"email": {Type: "string"},
			}},
			expectedType: "FIELD_ADDED",
			expectedPath: "email",
			breaking:     false,
		},
		{
			name:         "type changed",
			old:          models.Schema{Type: "string"},
			new:          models.Schema{Type: "integer"},
			expectedType: "TYPE_CHANGED",
			expectedPath: "",
			breaking:     true,
		},
		{
			name: "nested field type changed",
			old: models.Schema{Type: "object", Properties: map[string]models.Schema{
				"address": {Type: "object", Properties: map[string]models.Schema{
					"city": {Type: "string"},
				}},
			}},
			new: models.Schema{Type: "object", Properties: map[string]models.Schema{
				"address": {Type: "object", Properties: map[string]models.Schema{
					"city": {Type: "integer"},
				}},
			}},
			expectedType: "TYPE_CHANGED",
			expectedPath: "address.city",
			breaking:     true,
		},
		{
			name: "required removed",
			old: models.Schema{Type: "object", Required: []string{"name"},
				Properties: map[string]models.Schema{"name": {Type: "string"}}},
			new: models.Schema{Type: "object",
				Properties: map[string]models.Schema{"name": {Type: "string"}}},
			expectedType: "REQUIRED_REMOVED",
			expectedPath: "name",
			breaking:     true,
		},
		{
			name: "required added",
			old: models.Schema{Type: "object",
				Properties: map[string]models.Schema{"name": {Type: "string"}}},
			new: models.Schema{Type: "object", Required: []string{"name"},
				Properties: map[string]models.Schema{"name": {Type: "string"}}},
			expectedType: "REQUIRED_ADDED",
			expectedPath: "name",
			breaking:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldSpec := specWithSchema(tt.old)
			newSpec := specWithSchema(tt.new)
			result := CompareSpecs(oldSpec, newSpec)

			if len(result.Changes) != 1 {
				t.Fatalf("expected 1 change, got %d: %+v", len(result.Changes), result.Changes)
			}

			change := result.Changes[0]
			if change.Type != tt.expectedType || change.Field != tt.expectedPath || change.Breaking != tt.breaking {
				t.Errorf("unexpected change: %+v", change)
			}
		})
	}
}

func specWithSchema(schema models.Schema) models.OpenAPISpec {
	return models.OpenAPISpec{
		Paths: map[string]map[string]models.Operation{
			"/users/{id}": {
				"get": {
					Responses: map[string]models.Response{
						"200": {
							Content: map[string]models.MediaType{
								"application/json": {Schema: schema},
							},
						},
					},
				},
			},
		},
	}
}

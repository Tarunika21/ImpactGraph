package openapidiff

import (
	"testing"
)

func TestCompareSpecs(t *testing.T) {
	tests := []struct {
		name         string
		old          Schema
		new          Schema
		expectedType string
		expectedPath string
		breaking     bool
	}{
		{
			name: "field removed",
			old: Schema{Type: "object", Properties: map[string]Schema{
				"name": {Type: "string"},
			}},
			new:          Schema{Type: "object", Properties: map[string]Schema{}},
			expectedType: "FIELD_REMOVED",
			expectedPath: "name",
			breaking:     true,
		},
		{
			name: "field added",
			old:  Schema{Type: "object", Properties: map[string]Schema{}},
			new: Schema{Type: "object", Properties: map[string]Schema{
				"email": {Type: "string"},
			}},
			expectedType: "FIELD_ADDED",
			expectedPath: "email",
			breaking:     false,
		},
		{
			name:         "type changed",
			old:          Schema{Type: "string"},
			new:          Schema{Type: "integer"},
			expectedType: "TYPE_CHANGED",
			expectedPath: "",
			breaking:     true,
		},
		{
			name: "nested field type changed",
			old: Schema{Type: "object", Properties: map[string]Schema{
				"address": {Type: "object", Properties: map[string]Schema{
					"city": {Type: "string"},
				}},
			}},
			new: Schema{Type: "object", Properties: map[string]Schema{
				"address": {Type: "object", Properties: map[string]Schema{
					"city": {Type: "integer"},
				}},
			}},
			expectedType: "TYPE_CHANGED",
			expectedPath: "address.city",
			breaking:     true,
		},
		{
			name: "required removed",
			old: Schema{Type: "object", Required: []string{"name"},
				Properties: map[string]Schema{"name": {Type: "string"}}},
			new: Schema{Type: "object",
				Properties: map[string]Schema{"name": {Type: "string"}}},
			expectedType: "REQUIRED_REMOVED",
			expectedPath: "name",
			breaking:     true,
		},
		{
			name: "required added",
			old: Schema{Type: "object",
				Properties: map[string]Schema{"name": {Type: "string"}}},
			new: Schema{Type: "object", Required: []string{"name"},
				Properties: map[string]Schema{"name": {Type: "string"}}},
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

func specWithSchema(schema Schema) OpenAPISpec {
	return OpenAPISpec{
		Paths: map[string]map[string]Operation{
			"/users/{id}": {
				"get": {
					Responses: map[string]Response{
						"200": {
							Content: map[string]MediaType{
								"application/json": {Schema: schema},
							},
						},
					},
				},
			},
		},
	}
}

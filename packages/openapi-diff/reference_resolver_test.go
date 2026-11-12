package openapidiff

import (
	"strings"
	"testing"
)

func TestResolveSchema(t *testing.T) {
	tests := []struct {
		name       string
		schema     Schema
		components map[string]Schema
		wantField  string
		wantError  string
	}{
		{
			name:   "resolve direct reference",
			schema: Schema{Ref: "#/components/schemas/User"},
			components: map[string]Schema{
				"User": {
					Type: "object",
					Properties: map[string]Schema{
						"name": {Type: "string"},
					},
				},
			},
			wantField: "name",
		},
		{
			name:   "resolve nested reference",
			schema: Schema{Ref: "#/components/schemas/User"},
			components: map[string]Schema{
				"User": {
					Type: "object",
					Properties: map[string]Schema{
						"address": {Ref: "#/components/schemas/Address"},
					},
				},
				"Address": {
					Type: "object",
					Properties: map[string]Schema{
						"city": {Type: "string"},
					},
				},
			},
			wantField: "address",
		},
		{
			name:      "missing reference",
			schema:    Schema{Ref: "#/components/schemas/Missing"},
			wantError: "schema not found",
		},
		{
			name:   "circular reference",
			schema: Schema{Ref: "#/components/schemas/User"},
			components: map[string]Schema{
				"User": {
					Type: "object",
					Properties: map[string]Schema{
						"manager": {Ref: "#/components/schemas/User"},
					},
				},
			},
			wantError: "circular reference",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := OpenAPISpec{
				Components: Components{Schemas: tt.components},
			}

			resolved, err := ResolveSchema(tt.schema, spec)

			if tt.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantError) {
					t.Fatalf("expected error containing %q, got %v", tt.wantError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if _, exists := resolved.Properties[tt.wantField]; !exists {
				t.Errorf("expected property %q", tt.wantField)
			}
		})
	}
}

func TestCompareReferencedSchemas(t *testing.T) {
	oldSpec := specWithSchema(Schema{
		Ref: "#/components/schemas/User",
	})
	newSpec := specWithSchema(Schema{
		Ref: "#/components/schemas/User",
	})

	oldSpec.Components.Schemas = map[string]Schema{
		"User": {
			Type: "object",
			Properties: map[string]Schema{
				"id":   {Type: "integer"},
				"name": {Type: "string"},
			},
		},
	}

	newSpec.Components.Schemas = map[string]Schema{
		"User": {
			Type: "object",
			Properties: map[string]Schema{
				"id":   {Type: "string"},
				"name": {Type: "string"},
			},
		},
	}

	result, err := CompareSpecsChecked(oldSpec, newSpec)
	if err != nil {
		t.Fatal(err)
	}

	if result.TotalChanges != 1 || result.BreakingCount != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}

	change := result.Changes[0]
	if change.Type != "TYPE_CHANGED" || change.Field != "id" {
		t.Errorf("unexpected change: %+v", change)
	}
}

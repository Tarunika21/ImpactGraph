package services

import (
	"testing"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
)

func TestCompareSpecsDetectsRemovedField(t *testing.T) {
	oldSpec := makeSpec(map[string]models.Schema{
		"id":   {Type: "integer"},
		"name": {Type: "string"},
	})

	newSpec := makeSpec(map[string]models.Schema{
		"id": {Type: "integer"},
	})

	result := CompareSpecs(oldSpec, newSpec)

	if result.TotalChanges != 1 {
		t.Fatalf("expected 1 change, got %d", result.TotalChanges)
	}

	if result.BreakingCount != 1 {
		t.Fatalf("expected 1 breaking change, got %d", result.BreakingCount)
	}

	if result.Changes[0].Type != "FIELD_REMOVED" {
		t.Errorf("expected FIELD_REMOVED, got %s", result.Changes[0].Type)
	}
}

func TestCompareSpecsDetectsAddedField(t *testing.T) {
	oldSpec := makeSpec(map[string]models.Schema{
		"id": {Type: "integer"},
	})

	newSpec := makeSpec(map[string]models.Schema{
		"id":   {Type: "integer"},
		"name": {Type: "string"},
	})

	result := CompareSpecs(oldSpec, newSpec)

	if result.TotalChanges != 1 || result.BreakingCount != 0 {
		t.Fatalf("unexpected comparison result: %+v", result)
	}
}

func makeSpec(properties map[string]models.Schema) models.OpenAPISpec {
	return models.OpenAPISpec{
		Paths: map[string]map[string]models.Operation{
			"/users/{id}": {
				"get": {
					Responses: map[string]models.Response{
						"200": {
							Content: map[string]models.MediaType{
								"application/json": {
									Schema: models.Schema{
										Type:       "object",
										Properties: properties,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

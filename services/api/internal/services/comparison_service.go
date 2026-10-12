package services

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
)

func CompareSpecs(oldSpec, newSpec models.OpenAPISpec) models.ComparisonResponse {
	result, _ := CompareSpecsChecked(oldSpec, newSpec)
	return result
}

func CompareSpecsChecked(oldSpec, newSpec models.OpenAPISpec) (models.ComparisonResponse, error) {
	result := models.ComparisonResponse{Changes: []models.Change{}}

	for path, oldMethods := range oldSpec.Paths {
		newMethods := newSpec.Paths[path]

		for method, oldOperation := range oldMethods {
			endpoint := strings.ToUpper(method) + " " + path
			newOperation, exists := newMethods[method]

			if !exists {
				addChange(&result, models.Change{
					Type: "ENDPOINT_REMOVED", Endpoint: endpoint, Breaking: true,
				})
				continue
			}

			oldSchema, oldFound := getResponseSchema(oldOperation)
			newSchema, newFound := getResponseSchema(newOperation)

			if oldFound && !newFound {
				addChange(&result, models.Change{
					Type: "RESPONSE_REMOVED", Endpoint: endpoint, Breaking: true,
				})
				continue
			}
			if !oldFound || !newFound {
				continue
			}

			resolvedOld, err := ResolveSchema(oldSchema, oldSpec)
			if err != nil {
				return result, fmt.Errorf("%s old schema: %w", endpoint, err)
			}

			resolvedNew, err := ResolveSchema(newSchema, newSpec)
			if err != nil {
				return result, fmt.Errorf("%s new schema: %w", endpoint, err)
			}

			compareSchema(&result, endpoint, "", resolvedOld, resolvedNew)
		}
	}

	for path, newMethods := range newSpec.Paths {
		oldMethods := oldSpec.Paths[path]
		for method := range newMethods {
			if _, exists := oldMethods[method]; !exists {
				addChange(&result, models.Change{
					Type:     "ENDPOINT_ADDED",
					Endpoint: strings.ToUpper(method) + " " + path,
					Breaking: false,
				})
			}
		}
	}

	sort.Slice(result.Changes, func(i, j int) bool {
		a, b := result.Changes[i], result.Changes[j]
		if a.Endpoint != b.Endpoint {
			return a.Endpoint < b.Endpoint
		}
		if a.Field != b.Field {
			return a.Field < b.Field
		}
		return a.Type < b.Type
	})

	result.TotalChanges = len(result.Changes)
	return result, nil
}

func getResponseSchema(op models.Operation) (models.Schema, bool) {
	response, exists := op.Responses["200"]
	if !exists {
		return models.Schema{}, false
	}
	media, exists := response.Content["application/json"]
	if !exists {
		return models.Schema{}, false
	}
	return media.Schema, true
}

func compareSchema(result *models.ComparisonResponse, endpoint, prefix string, oldSchema, newSchema models.Schema) {
	if oldSchema.Type != newSchema.Type {
		addChange(result, models.Change{
			Type: "TYPE_CHANGED", Endpoint: endpoint, Field: prefix,
			OldValue: oldSchema.Type, NewValue: newSchema.Type, Breaking: true,
		})
		return
	}

	if oldSchema.Type == "array" && oldSchema.Items != nil && newSchema.Items != nil {
		compareSchema(result, endpoint, prefix+"[]", *oldSchema.Items, *newSchema.Items)
	}

	for field, oldProperty := range oldSchema.Properties {
		fieldPath := joinPath(prefix, field)
		newProperty, exists := newSchema.Properties[field]

		if !exists {
			addChange(result, models.Change{
				Type: "FIELD_REMOVED", Endpoint: endpoint,
				Field: fieldPath, Breaking: true,
			})
			continue
		}

		compareSchema(result, endpoint, fieldPath, oldProperty, newProperty)
	}

	for field := range newSchema.Properties {
		if _, exists := oldSchema.Properties[field]; !exists {
			addChange(result, models.Change{
				Type: "FIELD_ADDED", Endpoint: endpoint,
				Field: joinPath(prefix, field), Breaking: false,
			})
		}
	}

	oldRequired := toSet(oldSchema.Required)
	newRequired := toSet(newSchema.Required)

	for field := range oldRequired {
		if _, exists := newSchema.Properties[field]; !exists {
			continue
		}
		if !newRequired[field] {
			addChange(result, models.Change{
				Type: "REQUIRED_REMOVED", Endpoint: endpoint,
				Field: joinPath(prefix, field), Breaking: true,
			})
		}
	}

	for field := range newRequired {
		if _, exists := oldSchema.Properties[field]; !exists {
			continue
		}
		if !oldRequired[field] {
			addChange(result, models.Change{
				Type: "REQUIRED_ADDED", Endpoint: endpoint,
				Field: joinPath(prefix, field), Breaking: false,
			})
		}
	}
}

func joinPath(prefix, field string) string {
	if prefix == "" {
		return field
	}
	return prefix + "." + field
}

func toSet(values []string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
	}
	return result
}

func addChange(result *models.ComparisonResponse, change models.Change) {
	result.Changes = append(result.Changes, change)
	if change.Breaking {
		result.BreakingCount++
	}
}

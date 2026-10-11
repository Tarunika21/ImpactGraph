package services

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
)

func CompareSpecs(oldSpec, newSpec models.OpenAPISpec) models.ComparisonResponse {
	result := models.ComparisonResponse{
		Changes: []models.Change{},
	}

	for path, oldMethods := range oldSpec.Paths {
		newMethods, pathExists := newSpec.Paths[path]

		for method, oldOperation := range oldMethods {
			endpoint := strings.ToUpper(method) + " " + path

			newOperation, methodExists := newMethods[method]
			if !pathExists || !methodExists {
				result.Changes = append(result.Changes, models.Change{
					Type:     "ENDPOINT_REMOVED",
					Endpoint: endpoint,
					Breaking: true,
				})
				continue
			}

			oldFields := getResponseFields(oldOperation)
			newFields := getResponseFields(newOperation)

			for field := range oldFields {
				if _, exists := newFields[field]; !exists {
					result.Changes = append(result.Changes, models.Change{
						Type:     "FIELD_REMOVED",
						Endpoint: endpoint,
						Field:    field,
						Breaking: true,
					})
				}
			}

			for field := range newFields {
				if _, exists := oldFields[field]; !exists {
					result.Changes = append(result.Changes, models.Change{
						Type:     "FIELD_ADDED",
						Endpoint: endpoint,
						Field:    field,
						Breaking: false,
					})
				}
			}
		}
	}

	for path, newMethods := range newSpec.Paths {
		oldMethods := oldSpec.Paths[path]

		for method := range newMethods {
			if _, exists := oldMethods[method]; !exists {
				result.Changes = append(result.Changes, models.Change{
					Type:     "ENDPOINT_ADDED",
					Endpoint: strings.ToUpper(method) + " " + path,
					Breaking: false,
				})
			}
		}
	}

	sort.Slice(result.Changes, func(i, j int) bool {
		a, b := result.Changes[i], result.Changes[j]
		return fmt.Sprintf("%s:%s:%s", a.Endpoint, a.Type, a.Field) <
			fmt.Sprintf("%s:%s:%s", b.Endpoint, b.Type, b.Field)
	})

	result.TotalChanges = len(result.Changes)

	for _, change := range result.Changes {
		if change.Breaking {
			result.BreakingCount++
		}
	}

	return result
}

func getResponseFields(operation models.Operation) map[string]models.Schema {
	fields := make(map[string]models.Schema)

	response, exists := operation.Responses["200"]
	if !exists {
		return fields
	}

	content, exists := response.Content["application/json"]
	if !exists {
		return fields
	}

	for name, schema := range content.Schema.Properties {
		fields[name] = schema
	}

	return fields
}

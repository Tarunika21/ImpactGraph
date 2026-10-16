package services

import (
	"sort"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
)

func compareResponseStatusCodes(
	result *models.ComparisonResponse,
	endpoint string,
	oldResponses, newResponses map[string]models.Response,
) {
	for status := range oldResponses {
		if _, exists := newResponses[status]; !exists {
			addChange(result, models.Change{
				Type:     "RESPONSE_STATUS_REMOVED",
				Endpoint: endpoint,
				Field:    status,
				Breaking: true,
			})
		}
	}

	for status := range newResponses {
		if _, exists := oldResponses[status]; !exists {
			addChange(result, models.Change{
				Type:     "RESPONSE_STATUS_ADDED",
				Endpoint: endpoint,
				Field:    status,
				Breaking: false,
			})
		}
	}
}

func sortChanges(result *models.ComparisonResponse) {
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
}

package services

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
)

func AnalyzeImpact(
	changes []models.Change,
	consumers []models.ImpactConsumer,
) models.ImpactResponse {
	result := models.ImpactResponse{
		TotalConsumers:    len(consumers),
		ImpactedConsumers: []models.ImpactedConsumer{},
	}

	for _, consumer := range consumers {
		reasons := []string{}

		for _, change := range changes {
			if !change.Breaking {
				continue
			}

			if consumer.Endpoint != change.Endpoint {
				continue
			}

			if change.Field != "" &&
				!dependsOnField(consumer.Fields, change.Field) &&
				!affectsWholeEndpoint(change.Type) {
				continue
			}

			reasons = append(reasons, fmt.Sprintf(
				"%s: %s",
				change.Type,
				change.Field,
			))
		}

		if len(reasons) == 0 {
			continue
		}

		sort.Strings(reasons)

		result.ImpactedConsumers = append(
			result.ImpactedConsumers,
			models.ImpactedConsumer{
				ConsumerID:   consumer.ID,
				ConsumerName: consumer.Name,
				Endpoint:     consumer.Endpoint,
				Reasons:      reasons,
			},
		)
	}

	sort.Slice(result.ImpactedConsumers, func(i, j int) bool {
		return result.ImpactedConsumers[i].ConsumerID <
			result.ImpactedConsumers[j].ConsumerID
	})

	result.ImpactedCount = len(result.ImpactedConsumers)
	return result
}

func dependsOnField(fields []string, changedField string) bool {
	for _, field := range fields {
		if field == changedField ||
			strings.HasPrefix(field, changedField+".") ||
			strings.HasPrefix(changedField, field+".") {
			return true
		}
	}
	return false
}

func affectsWholeEndpoint(changeType string) bool {
	switch changeType {
	case "ENDPOINT_REMOVED", "RESPONSE_REMOVED",
		"RESPONSE_STATUS_REMOVED":
		return true
	default:
		return false
	}
}

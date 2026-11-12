package openapidiff

import (
	"fmt"
)

func compareParameters(
	result *ComparisonResponse,
	endpoint string,
	oldParams, newParams []Parameter,
) {
	oldMap := parameterMap(oldParams)
	newMap := parameterMap(newParams)

	for key, oldParam := range oldMap {
		newParam, exists := newMap[key]

		if !exists {
			addChange(result, Change{
				Type:     "PARAMETER_REMOVED",
				Endpoint: endpoint,
				Field:    key,
				Breaking: true,
			})
			continue
		}

		if oldParam.Schema.Type != newParam.Schema.Type {
			addChange(result, Change{
				Type:     "PARAMETER_TYPE_CHANGED",
				Endpoint: endpoint,
				Field:    key,
				OldValue: oldParam.Schema.Type,
				NewValue: newParam.Schema.Type,
				Breaking: true,
			})
		}

		if !oldParam.Required && newParam.Required {
			addChange(result, Change{
				Type:     "PARAMETER_BECAME_REQUIRED",
				Endpoint: endpoint,
				Field:    key,
				Breaking: true,
			})
		}

		if oldParam.Required && !newParam.Required {
			addChange(result, Change{
				Type:     "PARAMETER_BECAME_OPTIONAL",
				Endpoint: endpoint,
				Field:    key,
				Breaking: false,
			})
		}
	}

	for key, newParam := range newMap {
		if _, exists := oldMap[key]; exists {
			continue
		}

		addChange(result, Change{
			Type:     "PARAMETER_ADDED",
			Endpoint: endpoint,
			Field:    key,
			Breaking: newParam.Required,
		})
	}
}

func parameterMap(params []Parameter) map[string]Parameter {
	result := make(map[string]Parameter)

	for _, param := range params {
		key := fmt.Sprintf("%s:%s", param.In, param.Name)
		result[key] = param
	}

	return result
}

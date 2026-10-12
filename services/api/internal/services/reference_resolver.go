package services

import (
	"fmt"
	"strings"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
)

const schemaPrefix = "#/components/schemas/"

func ResolveSchema(
	schema models.Schema,
	spec models.OpenAPISpec,
) (models.Schema, error) {
	return resolveSchema(schema, spec, map[string]bool{}, 0)
}

func resolveSchema(
	schema models.Schema,
	spec models.OpenAPISpec,
	active map[string]bool,
	depth int,
) (models.Schema, error) {
	if depth > 100 {
		return models.Schema{}, fmt.Errorf("schema nesting limit exceeded")
	}

	if schema.Ref != "" {
		if !strings.HasPrefix(schema.Ref, schemaPrefix) {
			return models.Schema{}, fmt.Errorf("unsupported reference: %s", schema.Ref)
		}

		name := strings.TrimPrefix(schema.Ref, schemaPrefix)
		name = strings.ReplaceAll(name, "~1", "/")
		name = strings.ReplaceAll(name, "~0", "~")

		if active[name] {
			return models.Schema{}, fmt.Errorf("circular reference: %s", name)
		}

		target, exists := spec.Components.Schemas[name]
		if !exists {
			return models.Schema{}, fmt.Errorf("schema not found: %s", name)
		}

		active[name] = true
		resolved, err := resolveSchema(target, spec, active, depth+1)
		delete(active, name)
		return resolved, err
	}

	result := schema

	if schema.Properties != nil {
		result.Properties = make(map[string]models.Schema, len(schema.Properties))

		for name, property := range schema.Properties {
			resolved, err := resolveSchema(property, spec, active, depth+1)
			if err != nil {
				return models.Schema{}, fmt.Errorf("property %s: %w", name, err)
			}
			result.Properties[name] = resolved
		}
	}

	if schema.Items != nil {
		resolved, err := resolveSchema(*schema.Items, spec, active, depth+1)
		if err != nil {
			return models.Schema{}, fmt.Errorf("array items: %w", err)
		}
		result.Items = &resolved
	}

	return result, nil
}

package services

import (
	openapidiff "github.com/Tarunika21/ImpactGraph/packages/openapi-diff"
	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
)

func CompareSpecs(oldSpec, newSpec models.OpenAPISpec) models.ComparisonResponse {
	result, _ := CompareSpecsChecked(oldSpec, newSpec)
	return result
}

func CompareSpecsChecked(oldSpec, newSpec models.OpenAPISpec) (models.ComparisonResponse, error) {
	return openapidiff.Compare(oldSpec, newSpec)
}

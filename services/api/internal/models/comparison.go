package models

import (
	openapidiff "github.com/Tarunika21/ImpactGraph/packages/openapi-diff"
)

type OpenAPISpec = openapidiff.OpenAPISpec

type Components = openapidiff.Components

type Operation = openapidiff.Operation

type Parameter = openapidiff.Parameter

type Response = openapidiff.Response

type MediaType = openapidiff.MediaType

type Schema = openapidiff.Schema

type Change = openapidiff.Change

type ComparisonResponse = openapidiff.ComparisonResponse

type ComparisonRequest struct {
	OldSpec OpenAPISpec `json:"oldSpec"`
	NewSpec OpenAPISpec `json:"newSpec"`
}

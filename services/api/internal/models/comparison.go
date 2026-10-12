package models

type OpenAPISpec struct {
	Paths map[string]map[string]Operation `json:"paths"`
}

type Operation struct {
	Responses map[string]Response `json:"responses"`
}

type Response struct {
	Content map[string]MediaType `json:"content"`
}

type MediaType struct {
	Schema Schema `json:"schema"`
}

type Schema struct {
	Type       string            `json:"type"`
	Properties map[string]Schema `json:"properties"`
	Required   []string          `json:"required"`
	Items      *Schema           `json:"items,omitempty"`
}

type ComparisonRequest struct {
	OldSpec OpenAPISpec `json:"oldSpec"`
	NewSpec OpenAPISpec `json:"newSpec"`
}

type Change struct {
	Type     string `json:"type"`
	Endpoint string `json:"endpoint"`
	Field    string `json:"field,omitempty"`
	OldValue string `json:"oldValue,omitempty"`
	NewValue string `json:"newValue,omitempty"`
	Breaking bool   `json:"breaking"`
}

type ComparisonResponse struct {
	Changes       []Change `json:"changes"`
	TotalChanges  int      `json:"totalChanges"`
	BreakingCount int      `json:"breakingCount"`
}

package models

type OpenAPISpec struct {
	Paths      map[string]map[string]Operation `json:"paths"`
	Components Components                      `json:"components"`
}

type Components struct {
	Schemas map[string]Schema `json:"schemas"`
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
	Ref        string            `json:"$ref,omitempty"`
	Type       string            `json:"type,omitempty"`
	Properties map[string]Schema `json:"properties,omitempty"`
	Required   []string          `json:"required,omitempty"`
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

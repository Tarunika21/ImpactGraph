package services

import (
	"testing"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
)

func TestAnalyzeImpact(t *testing.T) {
	changes := []models.Change{
		{
			Type:     "FIELD_REMOVED",
			Endpoint: "GET /users/{id}",
			Field:    "name",
			Breaking: true,
		},
	}

	consumers := []models.ImpactConsumer{
		{
			ID:       "web",
			Name:     "Web App",
			Endpoint: "GET /users/{id}",
			Fields:   []string{"name", "email"},
		},
		{
			ID:       "mobile",
			Name:     "Mobile App",
			Endpoint: "GET /users/{id}",
			Fields:   []string{"email"},
		},
		{
			ID:       "admin",
			Name:     "Admin App",
			Endpoint: "GET /users/{id}",
			Fields:   []string{"name"},
		},
	}

	result := AnalyzeImpact(changes, consumers)

	if result.TotalConsumers != 3 {
		t.Errorf("expected 3 consumers, got %d", result.TotalConsumers)
	}

	if result.ImpactedCount != 2 {
		t.Errorf("expected 2 impacted consumers, got %d", result.ImpactedCount)
	}

	if len(result.ImpactedConsumers) != 2 {
		t.Fatalf("expected 2 impact records")
	}

	if result.ImpactedConsumers[0].ConsumerID != "admin" {
		t.Errorf("expected admin to be first")
	}

	if result.ImpactedConsumers[1].ConsumerID != "web" {
		t.Errorf("expected web to be second")
	}
}

func TestAnalyzeImpactEndpointRemoved(t *testing.T) {
	changes := []models.Change{
		{
			Type:     "ENDPOINT_REMOVED",
			Endpoint: "GET /users/{id}",
			Breaking: true,
		},
	}

	consumers := []models.ImpactConsumer{
		{
			ID:       "mobile",
			Name:     "Mobile App",
			Endpoint: "GET /users/{id}",
			Fields:   []string{"email"},
		},
	}

	result := AnalyzeImpact(changes, consumers)

	if result.ImpactedCount != 1 {
		t.Errorf("expected 1 impacted consumer")
	}
}

func TestAnalyzeImpactNonBreaking(t *testing.T) {
	changes := []models.Change{
		{
			Type:     "FIELD_ADDED",
			Endpoint: "GET /users/{id}",
			Field:    "age",
			Breaking: false,
		},
	}

	consumers := []models.ImpactConsumer{
		{
			ID:       "web",
			Name:     "Web App",
			Endpoint: "GET /users/{id}",
			Fields:   []string{"name"},
		},
	}

	result := AnalyzeImpact(changes, consumers)

	if result.ImpactedCount != 0 {
		t.Errorf("expected no impacted consumers")
	}
}

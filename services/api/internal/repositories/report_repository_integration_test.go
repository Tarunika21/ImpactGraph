package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestReportRepositoryIntegration(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not configured")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create database pool: %v", err)
	}
	t.Cleanup(func() {
		pool.Close()
	})

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("connect to test database: %v", err)
	}

	var databaseName string
	err = pool.QueryRow(ctx, "SELECT current_database()").Scan(&databaseName)
	if err != nil {
		t.Fatalf("check database name: %v", err)
	}
	if databaseName != "impactgraph_test" {
		t.Fatalf("refusing to run against %q", databaseName)
	}

	repository := NewReportRepository(pool)

	report := models.ImpactReport{
		ServiceName: "integration-test-service",
		OldVersion:  "1.0.0",
		NewVersion:  "2.0.0",
		Result: json.RawMessage(`{
			"comparison": {
				"totalChanges": 1,
				"breakingCount": 1
			},
			"impact": {
				"totalConsumers": 1,
				"impactedCount": 1
			}
		}`),
	}

	saved, err := repository.Create(ctx, report)
	if err != nil {
		t.Fatalf("save report: %v", err)
	}

	t.Cleanup(func() {
		_, err := pool.Exec(
			context.Background(),
			"DELETE FROM impact_reports WHERE id = $1",
			saved.ID,
		)
		if err != nil {
			t.Errorf("cleanup failed: %v", err)
		}
	})

	if saved.ID <= 0 {
		t.Fatalf("expected a generated ID, got %d", saved.ID)
	}

	retrieved, err := repository.GetByID(ctx, saved.ID)
	if err != nil {
		t.Fatalf("retrieve report: %v", err)
	}

	if retrieved.ServiceName != report.ServiceName {
		t.Errorf("service name mismatch")
	}

	if retrieved.OldVersion != report.OldVersion {
		t.Errorf("old version mismatch")
	}

	if retrieved.NewVersion != report.NewVersion {
		t.Errorf("new version mismatch")
	}

	var result map[string]any
	if err := json.Unmarshal(retrieved.Result, &result); err != nil {
		t.Fatalf("invalid stored JSON: %v", err)
	}

	if _, ok := result["comparison"]; !ok {
		t.Error("comparison missing from stored report")
	}

	if _, ok := result["impact"]; !ok {
		t.Error("impact missing from stored report")
	}

	reports, err := repository.GetByService(ctx, report.ServiceName)
	if err != nil {
		t.Fatalf("filter reports: %v", err)
	}

	found := false
	for _, item := range reports {
		if item.ID == saved.ID {
			found = true
			break
		}
	}

	if !found {
		t.Error("saved report not found in filtered results")
	}

	_, err = repository.GetByID(ctx, -1)
	if !errors.Is(err, ErrReportNotFound) {
		t.Errorf("expected ErrReportNotFound, got %v", err)
	}
}

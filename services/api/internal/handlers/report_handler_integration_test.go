package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
	"github.com/Tarunika21/ImpactGraph/services/api/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestReportHandlerIntegration(t *testing.T) {
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
	err = pool.QueryRow(
		ctx,
		"SELECT current_database()",
	).Scan(&databaseName)
	if err != nil {
		t.Fatalf("check database name: %v", err)
	}

	if databaseName != "impactgraph_test" {
		t.Fatalf(
			"refusing to run against database %q",
			databaseName,
		)
	}

	repository := repositories.NewReportRepository(pool)
	handler := NewReportHandler(repository)

	report := models.ImpactReport{
		ServiceName: "http-integration-test-service",
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
		t.Fatalf("create test report: %v", err)
	}

	t.Cleanup(func() {
		_, err := pool.Exec(
			context.Background(),
			"DELETE FROM impact_reports WHERE id = $1",
			saved.ID,
		)
		if err != nil {
			t.Errorf("cleanup test report: %v", err)
		}
	})

	t.Run("Get report by ID", func(t *testing.T) {
		request := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/reports/%d", saved.ID),
			nil,
		)
		recorder := httptest.NewRecorder()

		handler.HandleReportByID(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf(
				"expected HTTP 200, got %d: %s",
				recorder.Code,
				recorder.Body.String(),
			)
		}

		var response models.ImpactReport
		if err := json.Unmarshal(
			recorder.Body.Bytes(),
			&response,
		); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.ID != saved.ID {
			t.Errorf(
				"expected report ID %d, got %d",
				saved.ID,
				response.ID,
			)
		}

		if response.ServiceName != report.ServiceName {
			t.Errorf("service name mismatch")
		}
	})

	t.Run("Filter reports by service", func(t *testing.T) {
		request := httptest.NewRequest(
			http.MethodGet,
			"/api/v1/reports?serviceName=http-integration-test-service",
			nil,
		)
		recorder := httptest.NewRecorder()

		handler.HandleReports(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf(
				"expected HTTP 200, got %d: %s",
				recorder.Code,
				recorder.Body.String(),
			)
		}

		var reports []models.ImpactReport
		if err := json.Unmarshal(
			recorder.Body.Bytes(),
			&reports,
		); err != nil {
			t.Fatalf("decode reports: %v", err)
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
	})

	t.Run("Get all reports", func(t *testing.T) {
		request := httptest.NewRequest(
			http.MethodGet,
			"/api/v1/reports",
			nil,
		)
		recorder := httptest.NewRecorder()

		handler.HandleReports(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf(
				"expected HTTP 200, got %d: %s",
				recorder.Code,
				recorder.Body.String(),
			)
		}

		var reports []models.ImpactReport
		if err := json.Unmarshal(
			recorder.Body.Bytes(),
			&reports,
		); err != nil {
			t.Fatalf("decode reports: %v", err)
		}

		found := false
		for _, item := range reports {
			if item.ID == saved.ID {
				found = true
				break
			}
		}

		if !found {
			t.Error("saved report not found in report history")
		}
	})

	t.Run("Nonexistent report returns 404", func(t *testing.T) {
		request := httptest.NewRequest(
			http.MethodGet,
			"/api/v1/reports/9223372036854775807",
			nil,
		)
		recorder := httptest.NewRecorder()

		handler.HandleReportByID(recorder, request)

		if recorder.Code != http.StatusNotFound {
			t.Errorf(
				"expected HTTP 404, got %d",
				recorder.Code,
			)
		}
	})
}

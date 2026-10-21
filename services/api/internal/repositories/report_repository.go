package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrReportNotFound = errors.New("report not found")

type ReportRepository struct {
	db *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Create(
	ctx context.Context,
	report models.ImpactReport,
) (models.ImpactReport, error) {
	query := `
		INSERT INTO impact_reports
			(service_name, old_version, new_version, result)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		report.ServiceName,
		report.OldVersion,
		report.NewVersion,
		report.Result,
	).Scan(&report.ID, &report.CreatedAt)

	if err != nil {
		return models.ImpactReport{}, fmt.Errorf("create impact report: %w", err)
	}

	return report, nil
}

func (r *ReportRepository) GetAll(
	ctx context.Context,
) ([]models.ImpactReport, error) {
	query := `
		SELECT id, service_name, old_version, new_version, result, created_at
		FROM impact_reports
		ORDER BY created_at DESC, id DESC
	`

	return r.queryReports(ctx, query)
}

func (r *ReportRepository) GetByService(
	ctx context.Context,
	serviceName string,
) ([]models.ImpactReport, error) {
	query := `
		SELECT id, service_name, old_version, new_version, result, created_at
		FROM impact_reports
		WHERE service_name = $1
		ORDER BY created_at DESC, id DESC
	`

	return r.queryReports(ctx, query, serviceName)
}

func (r *ReportRepository) GetByID(
	ctx context.Context,
	id int64,
) (models.ImpactReport, error) {
	query := `
		SELECT id, service_name, old_version, new_version, result, created_at
		FROM impact_reports
		WHERE id = $1
	`

	var report models.ImpactReport

	err := r.db.QueryRow(ctx, query, id).Scan(
		&report.ID,
		&report.ServiceName,
		&report.OldVersion,
		&report.NewVersion,
		&report.Result,
		&report.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.ImpactReport{}, ErrReportNotFound
	}
	if err != nil {
		return models.ImpactReport{}, fmt.Errorf("get report by ID: %w", err)
	}

	return report, nil
}

func (r *ReportRepository) queryReports(
	ctx context.Context,
	query string,
	args ...any,
) ([]models.ImpactReport, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query impact reports: %w", err)
	}
	defer rows.Close()

	reports := []models.ImpactReport{}

	for rows.Next() {
		var report models.ImpactReport

		err := rows.Scan(
			&report.ID,
			&report.ServiceName,
			&report.OldVersion,
			&report.NewVersion,
			&report.Result,
			&report.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan impact report: %w", err)
		}

		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate impact reports: %w", err)
	}

	return reports, nil
}

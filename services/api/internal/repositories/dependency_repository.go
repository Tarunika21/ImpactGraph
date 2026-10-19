package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrDependencyExists = errors.New("dependency already exists")
var ErrConsumerNotFound = errors.New("consumer not found")

type DependencyRepository struct {
	db *pgxpool.Pool
}

func NewDependencyRepository(db *pgxpool.Pool) *DependencyRepository {
	return &DependencyRepository{db: db}
}

func (r *DependencyRepository) Create(
	ctx context.Context,
	request models.CreateConsumerDependencyRequest,
) (models.ConsumerDependency, error) {
	var dependency models.ConsumerDependency

	query := `
		INSERT INTO consumer_dependencies
			(consumer_id, service_name, endpoint, fields)
		VALUES ($1, $2, $3, $4)
		RETURNING id, consumer_id, service_name, endpoint, fields, created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		request.ConsumerID,
		request.ServiceName,
		request.Endpoint,
		request.Fields,
	).Scan(
		&dependency.ID,
		&dependency.ConsumerID,
		&dependency.ServiceName,
		&dependency.Endpoint,
		&dependency.Fields,
		&dependency.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return dependency, ErrDependencyExists
			case "23503":
				return dependency, ErrConsumerNotFound
			}
		}
		return dependency, fmt.Errorf("create dependency: %w", err)
	}

	return dependency, nil
}

func (r *DependencyRepository) GetAll(
	ctx context.Context,
) ([]models.ConsumerDependency, error) {
	query := `
		SELECT id, consumer_id, service_name, endpoint, fields, created_at
		FROM consumer_dependencies
		ORDER BY created_at DESC, id DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query dependencies: %w", err)
	}
	defer rows.Close()

	dependencies := []models.ConsumerDependency{}

	for rows.Next() {
		var dependency models.ConsumerDependency

		if err := rows.Scan(
			&dependency.ID,
			&dependency.ConsumerID,
			&dependency.ServiceName,
			&dependency.Endpoint,
			&dependency.Fields,
			&dependency.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan dependency: %w", err)
		}

		dependencies = append(dependencies, dependency)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate dependencies: %w", err)
	}

	return dependencies, nil
}

func (r *DependencyRepository) GetImpactConsumers(
	ctx context.Context,
	serviceName string,
) ([]models.ImpactConsumer, error) {
	query := `
		SELECT c.id, c.name, d.endpoint, d.fields
		FROM consumer_dependencies d
		INNER JOIN consumers c ON c.id = d.consumer_id
		WHERE d.service_name = $1
		ORDER BY c.id, d.endpoint
	`

	rows, err := r.db.Query(ctx, query, serviceName)
	if err != nil {
		return nil, fmt.Errorf("query impact consumers: %w", err)
	}
	defer rows.Close()

	consumers := []models.ImpactConsumer{}

	for rows.Next() {
		var consumer models.ImpactConsumer
		err := rows.Scan(
			&consumer.ID,
			&consumer.Name,
			&consumer.Endpoint,
			&consumer.Fields,
		)
		if err != nil {
			return nil, fmt.Errorf("scan impact consumer: %w", err)
		}
		consumers = append(consumers, consumer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate impact consumers: %w", err)
	}

	return consumers, nil
}

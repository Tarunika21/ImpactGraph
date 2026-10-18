package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrConsumerExists = errors.New("consumer already exists")

type ConsumerRepository struct {
	db *pgxpool.Pool
}

func NewConsumerRepository(db *pgxpool.Pool) *ConsumerRepository {
	return &ConsumerRepository{db: db}
}

func (r *ConsumerRepository) Create(
	ctx context.Context,
	request models.CreateConsumerRequest,
) (models.Consumer, error) {
	var consumer models.Consumer

	query := `
		INSERT INTO consumers (id, name)
		VALUES ($1, $2)
		RETURNING id, name, created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		request.ID,
		request.Name,
	).Scan(
		&consumer.ID,
		&consumer.Name,
		&consumer.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return consumer, ErrConsumerExists
		}
		return consumer, fmt.Errorf("create consumer: %w", err)
	}

	return consumer, nil
}

func (r *ConsumerRepository) GetAll(
	ctx context.Context,
) ([]models.Consumer, error) {
	query := `
		SELECT id, name, created_at
		FROM consumers
		ORDER BY created_at DESC, id DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query consumers: %w", err)
	}
	defer rows.Close()

	consumers := []models.Consumer{}

	for rows.Next() {
		var consumer models.Consumer

		if err := rows.Scan(
			&consumer.ID,
			&consumer.Name,
			&consumer.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan consumer: %w", err)
		}

		consumers = append(consumers, consumer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate consumers: %w", err)
	}

	return consumers, nil
}

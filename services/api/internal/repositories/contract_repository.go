package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrContractExists = errors.New("contract version already exists")

type ContractRepository struct {
	db *pgxpool.Pool
}

func NewContractRepository(db *pgxpool.Pool) *ContractRepository {
	return &ContractRepository{db: db}
}

func (r *ContractRepository) Create(
	ctx context.Context,
	request models.CreateContractRequest,
) (models.Contract, error) {
	var contract models.Contract

	query := `
		INSERT INTO api_contracts (service_name, version, specification)
		VALUES ($1, $2, $3)
		RETURNING id, service_name, version, specification, created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		request.ServiceName,
		request.Version,
		request.Specification,
	).Scan(
		&contract.ID,
		&contract.ServiceName,
		&contract.Version,
		&contract.Specification,
		&contract.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return contract, ErrContractExists
		}
		return contract, fmt.Errorf("create contract: %w", err)
	}

	return contract, nil
}

func (r *ContractRepository) GetAll(
	ctx context.Context,
) ([]models.Contract, error) {
	query := `
		SELECT id, service_name, version, specification, created_at
		FROM api_contracts
		ORDER BY created_at DESC, id DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query contracts: %w", err)
	}
	defer rows.Close()

	contracts := []models.Contract{}

	for rows.Next() {
		var contract models.Contract

		err := rows.Scan(
			&contract.ID,
			&contract.ServiceName,
			&contract.Version,
			&contract.Specification,
			&contract.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan contract: %w", err)
		}

		contracts = append(contracts, contract)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate contracts: %w", err)
	}

	return contracts, nil
}

var _ = pgx.ErrNoRows

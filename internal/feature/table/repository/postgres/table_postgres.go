package postgres

import (
	"Test/internal/feature/table/entity"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTableRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresTableRepository(pool *pgxpool.Pool) *PostgresTableRepository {
	return &PostgresTableRepository{pool: pool}
}

func (r *PostgresTableRepository) Create(ctx context.Context, table *entity.Table) error {
	query := `INSERT INTO tables (organization_id, number, seats) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, query, table.OrganizationID, table.Number, table.Seats)
	return err
}

func (r *PostgresTableRepository) UpdateSeats(ctx context.Context, seats uint64, id uuid.UUID, organizationID uuid.UUID) error {
	query := `UPDATE tables SET seats = $3 WHERE id = $1 AND organization_id = $2`
	_, err := r.pool.Exec(ctx, query, id, organizationID, seats)
	return err
}

func (r *PostgresTableRepository) UpdateNumber(ctx context.Context, number uint64, id uuid.UUID, organizationID uuid.UUID) error {
	query := `UPDATE tables SET number = $3 WHERE id = $1 AND organization_id = $2`
	_, err := r.pool.Exec(ctx, query, id, organizationID, number)
	return err
}

func (r *PostgresTableRepository) UpdateStatus(ctx context.Context, status string, id uuid.UUID, organizationID uuid.UUID) error {
	query := `UPDATE tables SET status = $3 WHERE id = $1 AND organization_id = $2`
	_, err := r.pool.Exec(ctx, query, id, organizationID, status)
	return err
}

func (r *PostgresTableRepository) Get(ctx context.Context, id uuid.UUID) (*entity.Table, error) {
	query := `SELECT organization_id, number, status, seats, created_at, updated_at FROM tables WHERE id = $1`
	var table *entity.Table
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&table.OrganizationID,
		&table.Number,
		&table.Status,
		&table.Seats,
		&table.CreatedAt,
		&table.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return table, err
}

func (r *PostgresTableRepository) GetAll(ctx context.Context, organizationID uuid.UUID) ([]*entity.Table, error) {
	query := `SELECT id, number, status, seats, created_at, updated_at FROM tables WHERE organization_id = $1`
	rows, err := r.pool.Query(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []*entity.Table
	for rows.Next() {
		table := &entity.Table{
			OrganizationID: organizationID,
		}
		err = rows.Scan(
			&table.ID,
			&table.Number,
			&table.Status,
			&table.Seats,
			&table.CreatedAt,
			&table.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return tables, nil
}

func (r *PostgresTableRepository) Delete(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) error {
	query := `DELETE FROM tables WHERE id = $1 AND organization_id = $2 `
	_, err := r.pool.Exec(ctx, query, id, organizationID)
	return err
}

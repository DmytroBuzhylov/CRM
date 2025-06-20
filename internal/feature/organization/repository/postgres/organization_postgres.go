package postgres

import (
	"Test/internal/feature/organization/entity"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PostgresOrganizationRepository struct {
	db *pgxpool.Pool
}

func NewPostgresOrganizationRepository(db *pgxpool.Pool) *PostgresOrganizationRepository {
	return &PostgresOrganizationRepository{db: db}
}

func (r *PostgresOrganizationRepository) Create(ctx context.Context, organization entity.Organization) error {
	query := `INSERT INTO organizations (name, description, owner_user_id, users, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, organization.Name, organization.Description, organization.OwnerUserID, organization.Users, time.Now())
	if err != nil {
		return err
	}
	return nil
}

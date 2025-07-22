package postgres

import (
	"Test/internal/feature/organization/entity"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PostgresOrganizationRepository struct {
	db *pgxpool.Pool
}

func NewPostgresOrganizationRepository(db *pgxpool.Pool) *PostgresOrganizationRepository {
	return &PostgresOrganizationRepository{db: db}
}

func (r *PostgresOrganizationRepository) Create(ctx context.Context, organization *entity.Organization) error {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`, organization.OwnerUserID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("user with id %d does not exist", organization.OwnerUserID)
	}

	queryToCreateOrganization := `INSERT INTO organizations (name, description, owner_user_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	queryToCreateOrganizations_Users := `INSERT INTO organizations_users (organization_id, user_id, role) VALUES ($1, $2, $3)`

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var organizationID uint64
	err = tx.QueryRow(ctx, queryToCreateOrganization, organization.Name, organization.Description, organization.OwnerUserID, time.Now()).Scan(&organizationID)
	if err != nil {
		return fmt.Errorf("failed to create organization: %w", err)
	}

	_, err = tx.Exec(ctx, queryToCreateOrganizations_Users, organizationID, organization.OwnerUserID, "owner")
	if err != nil {
		return fmt.Errorf("failed to create organization-user link: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	organization.ID = organizationID
	return nil
}

package postgres

import (
	"Test/internal/feature/organization/entity"
	entity2 "Test/internal/feature/user/entity"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresInvitationRepository struct {
	db *pgxpool.Pool
}

func NewPostgresInvitationRepository(db *pgxpool.Pool) *PostgresInvitationRepository {
	return &PostgresInvitationRepository{db: db}
}

func (r *PostgresInvitationRepository) Save(ctx context.Context, inv entity.Invitation) error {
	query := `INSERT INTO invitations (organization_id, invited_email, invitation_code, expires_at, status, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.db.QueryRow(ctx, query, inv.OrganizationID, inv.InvitedEmail, inv.InvitationCode, inv.ExpiresAt, inv.Status, inv.CreatedAt).Scan(&inv.ID)
	return err
}

func (r *PostgresInvitationRepository) GetByCode(ctx context.Context, code string) (entity.Invitation, error) {
	var inv entity.Invitation
	query := `SELECT id, organization_id, invited_email, invitation_code, expires_at, status, created_at FROM invitations WHERE code = $1`
	err := r.db.QueryRow(ctx, query, code).Scan(&inv.ID, inv.OrganizationID, inv.InvitedEmail, inv.InvitationCode, inv.ExpiresAt, inv.Status, inv.CreatedAt)
	return inv, err
}

func (r *PostgresInvitationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE invitation SET status = $2 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, status)
	return err
}

func (r *PostgresInvitationRepository) AcceptInvite(ctx context.Context, userID uuid.UUID, code string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer tx.Rollback(ctx)

	var inv entity.Invitation
	var user entity2.User
	query := `SELECT
             i.id AS invitation_id,
             i.organization_id,
             i.invited_email,
             i.invitation_code,
             i.expires_at,
             i.status,
             u.id AS user_id,
             u.email AS user_email
             FROM
                 invitations i
             INNER JOIN 
                users u ON i.invited_email = u.email
             WHERE 
                 i.invitation_code = $1
             AND u.id = $2
             AND i.status = 'pending'
             AND i.expires_at > now();
    `
	err = tx.QueryRow(ctx, query, code, userID).Scan(&inv.ID, &inv.OrganizationID, &inv.InvitedEmail, &inv.InvitationCode, &inv.ExpiresAt, &inv.Status, &user.ID, &user.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("приглашение не найдено, истекло или недействительно для этого пользователя")
		}
		return fmt.Errorf("не удалось получить детали приглашения: %w", err)
	}

	var exists bool
	err = tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM organizations_users WHERE user_id = $1)`, user.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("не удалось проверить, состоит ли пользователь в другой организации: %w", err)
	}
	if exists {
		return errors.New("пользователь уже состоит в другой организации")
	}

	_, err = tx.Exec(ctx, `INSERT INTO organizations_users (organization_id, user_id, role) VALUES ($1, $2, $3)`, inv.OrganizationID, user.ID, "member")
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errors.New("пользователь уже является членом этой организации")
		}
		return fmt.Errorf("не удалось добавить пользователя в организацию: %w", err)
	}

	_, err = tx.Exec(ctx, `UPDATE invitations SET status = $1 WHERE id = $2`, "accepted", inv.ID)
	if err != nil {
		return fmt.Errorf("не удалось обновить статус приглашения: %w", err)
	}

	return tx.Commit(ctx)
}

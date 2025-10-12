package postgres

import (
	"Test/internal/feature/notification/entity"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresNotificationRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresNotificationRepository(pool *pgxpool.Pool) *postgresNotificationRepository {
	return &postgresNotificationRepository{pool: pool}
}

func (r *postgresNotificationRepository) Create(ctx context.Context, notify *entity.Notification) error {
	query := `
		INSERT INTO notifications (
			id, organization_id, user_id, type, text,
			send_to_email, send_to_app, viewed,
			resource_type, related_resource_id, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id
	`
	row := r.pool.QueryRow(ctx, query,
		notify.ID,
		notify.OrganizationID,
		notify.UserID,
		notify.Type,
		notify.Text,
		notify.SendToEmail,
		notify.SendToApp,
		notify.Viewed,
		func() any {
			if notify.ResourceType != nil {
				return notify.ResourceType
			}
			return nil
		}(),
		notify.RelatedResourceID,
		notify.CreatedAt,
	)

	return row.Scan(&notify.ID)
}

func (r *postgresNotificationRepository) UpdateViewStatus(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE notifications SET viewed = true WHERE id = $1`
	cmd, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("notification not found")
	}
	return nil
}

func (r *postgresNotificationRepository) Get(ctx context.Context, id uuid.UUID) (*entity.Notification, error) {
	query := `
		SELECT id, organization_id, user_id, type, text,
		       send_to_email, send_to_app, viewed,
		       resource_type, related_resource_id, created_at
		FROM notifications
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)

	var n entity.Notification
	//var typeStr, resTypeStr *string

	err := row.Scan(
		&n.ID,
		&n.OrganizationID,
		&n.UserID,
		&n.Type,
		&n.Text,
		&n.SendToEmail,
		&n.SendToApp,
		&n.Viewed,
		&n.ResourceType,
		&n.RelatedResourceID,
		&n.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *postgresNotificationRepository) GetMany(ctx context.Context, organizationID uuid.UUID) ([]*entity.Notification, error) {
	query := `
		SELECT id, organization_id, user_id, type, text,
		       send_to_email, send_to_app, viewed,
		       resource_type, related_resource_id, created_at
		FROM notifications
		WHERE organization_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.Notification

	for rows.Next() {
		var n entity.Notification

		err := rows.Scan(
			&n.ID,
			&n.OrganizationID,
			&n.UserID,
			&n.Type,
			&n.Text,
			&n.SendToEmail,
			&n.SendToApp,
			&n.Viewed,
			&n.ResourceType,
			&n.RelatedResourceID,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, &n)
	}

	return result, nil
}

func (r *postgresNotificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM notifications WHERE id = $1`
	cmd, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("notification not found")
	}
	return nil
}

func (r *postgresNotificationRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error) {
	query := `
		SELECT id, organization_id, user_id, type, text,
		       send_to_email, send_to_app, viewed,
		       resource_type, related_resource_id, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.Notification
	for rows.Next() {
		var n entity.Notification
		if err := rows.Scan(
			&n.ID,
			&n.OrganizationID,
			&n.UserID,
			&n.Type,
			&n.Text,
			&n.SendToEmail,
			&n.SendToApp,
			&n.Viewed,
			&n.ResourceType,
			&n.RelatedResourceID,
			&n.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, &n)
		fmt.Println(n.Viewed)
	}
	return result, nil
}

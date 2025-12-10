package repository

import (
	"Test/internal/feature/permission/entity"
	"Test/internal/pkg/provider"
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type RolePermissionRepo struct {
	connector provider.SQLConnector
	qb        sq.StatementBuilderType
}

func NewRolePermissionRepo(connector provider.SQLConnector) *RolePermissionRepo {
	return &RolePermissionRepo{
		connector: connector,
		qb:        sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *RolePermissionRepo) AddPermissionToRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	insertBuilder := r.qb.
		Insert("role_permissions").
		Columns("role_id", "permission_id", "created_at")

	now := time.Now()

	for _, permissionID := range permissionIDs {
		insertBuilder = insertBuilder.Values(roleID, permissionID, now)
	}

	query, args, err := insertBuilder.
		Suffix("ON CONFLICT (role_id, permission_id) DO NOTHING").
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.connector.Exec(ctx, query, args...)
	return err
}

func (r *RolePermissionRepo) RemovePermissionFromRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	if len(permissionIDs) == 0 {
		return nil
	}
	query, args, err := r.qb.
		Delete("role_permissions").
		Where(sq.Eq{"role_id": roleID, "permission_id": permissionIDs}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.connector.Exec(ctx, query, args...)
	return err
}

func (r *RolePermissionRepo) GetPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]*entity.Permission, error) {
	query, args, err := r.qb.
		Select("p.id", "p.name", "p.description", "p.created_at", "p.updated_at").
		From("permissions p").
		Join("role_permissions rp ON rp.permission_id = p.id").
		Where(sq.Eq{"rp.role_id": roleID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.connector.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []*entity.Permission
	for rows.Next() {
		var p entity.Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		perms = append(perms, &p)
	}

	return perms, nil
}

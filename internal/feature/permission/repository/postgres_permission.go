package repository

import (
	"Test/internal/feature/permission/entity"
	"Test/internal/pkg/provider"
	"context"
	_ "embed"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type PermissionRepo struct {
	provider provider.DataProvider
	qb       sq.StatementBuilderType
}

func NewPermissionRepo(provider provider.DataProvider) *PermissionRepo {
	return &PermissionRepo{
		provider: provider,
		qb:       sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PermissionRepo) Create(ctx context.Context, perm *entity.Permission) error {
	query, args, err := r.qb.
		Insert("permissions").
		Columns("id", "name", "description", "created_at", "updated_at").
		Values(perm.ID, perm.Name, perm.Description, time.Now(), time.Now()).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.provider.Exec(ctx, query, args...)
	return err
}

func (r *PermissionRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error) {
	query, args, err := r.qb.
		Select("id", "name", "description", "created_at", "updated_at").
		From("permissions").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.provider.QueryRow(ctx, query, args...)

	var p entity.Permission
	err = row.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PermissionRepo) GetAll(ctx context.Context) ([]*entity.Permission, error) {
	query, args, err := r.qb.
		Select("id", "name", "description", "created_at", "updated_at").
		From("permissions").
		OrderBy("name ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.provider.Query(ctx, query, args...)
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

//go:embed sql/permission/permission_exists.sql
var permissionExistsQuery string

func (r *PermissionRepo) PermissionExists(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID, permissionName string) (bool, error) {
	var hasPermission bool
	err := r.provider.QueryRow(ctx, permissionExistsQuery, userID, organizationID, permissionName).Scan(&hasPermission)
	if err != nil {
		return false, err
	}
	return hasPermission, nil
}

func (r *PermissionRepo) GetByName(ctx context.Context, name string) (*entity.Permission, error) {
	query, args, err := r.qb.
		Select("id", "name", "description", "created_at", "updated_at").
		From("permissions").
		Where(sq.Eq{"name": name}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.provider.QueryRow(ctx, query, args...)

	var p entity.Permission
	err = row.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PermissionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := r.qb.
		Delete("permissions").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.provider.Exec(ctx, query, args...)
	return err
}

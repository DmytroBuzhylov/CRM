package repository

import (
	"Test/internal/feature/permission/entity"
	"Test/internal/pkg/provider"
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type RoleRepo struct {
	provider provider.DataProvider
	qb       sq.StatementBuilderType
}

func NewRoleRepo(provider provider.DataProvider) *RoleRepo {
	return &RoleRepo{
		provider: provider,
		qb:       sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *RoleRepo) Create(ctx context.Context, role *entity.Role) error {
	query, args, err := r.qb.
		Insert("roles").
		Columns("id", "organization_id", "name", "description", "created_at", "updated_at").
		Values(role.ID, role.OrganizationID, role.Name, role.Description, time.Now(), time.Now()).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.provider.Exec(ctx, query, args...)
	return err
}

func (r *RoleRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	query, args, err := r.qb.
		Select("id", "organization_id", "name", "description", "created_at", "updated_at").
		From("roles").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.provider.QueryRow(ctx, query, args...)

	var role entity.Role
	err = row.Scan(&role.ID, &role.OrganizationID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepo) Update(ctx context.Context, role *entity.Role) error {
	query, args, err := r.qb.
		Update("roles").
		Set("name", role.Name).
		Set("description", role.Description).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": role.ID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.provider.Exec(ctx, query, args...)
	return err
}

func (r *RoleRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := r.qb.
		Delete("roles").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.provider.Exec(ctx, query, args...)
	return err
}

func (r *RoleRepo) GetMany(ctx context.Context, organizationID uuid.UUID) ([]*entity.Role, error) {
	query, args, err := r.qb.
		Select("id", "organization_id", "name", "description", "created_at", "updated_at").
		From("roles").
		Where(sq.Eq{"organization_id": organizationID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.provider.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entity.Role
	for rows.Next() {
		var role entity.Role
		if err := rows.Scan(&role.ID, &role.OrganizationID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	return roles, nil
}

package postgres

import (
	"Test/internal/feature/inventory/entity"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PostgresInventoryRepository struct {
	db *pgxpool.Pool
}

func NewPostgresInventoryRepository(db *pgxpool.Pool) *PostgresInventoryRepository {
	return &PostgresInventoryRepository{db: db}
}

func (r *PostgresInventoryRepository) Get(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) (entity.Ingredient, error) {
	var ingredient entity.Ingredient
	ingredient.ID, ingredient.OrganizationID = id, organizationID

	query := `SELECT name, unit, quantity, minimum_quantity, created_at, updated_at FROM ingredients WHERE id = $1 AND organization_id = $2`
	err := r.db.QueryRow(ctx, query, ingredient.ID, ingredient.OrganizationID).Scan(&ingredient.Name, &ingredient.Unit, &ingredient.Quantity, &ingredient.MinimumQuantity, &ingredient.CreatedAt, &ingredient.UpdatedAt)
	if err != nil {
		return entity.Ingredient{}, err
	}
	return ingredient, nil
}

func (r *PostgresInventoryRepository) GetByName(ctx context.Context, name string) (entity.Ingredient, error) {
	var ingredient entity.Ingredient
	query := `SELECT id, organization_id, name, unit, quantity, minimum_quantity, created_at, updated_at FROM ingredients WHERE name = $1`
	err := r.db.QueryRow(ctx, query, fmt.Sprintf("%"+name+"%")).Scan(&ingredient.ID, &ingredient.OrganizationID, &ingredient.Name, &ingredient.Unit, &ingredient.Quantity, &ingredient.MinimumQuantity, &ingredient.CreatedAt, &ingredient.UpdatedAt)
	if err != nil {
		return entity.Ingredient{}, err
	}
	return ingredient, nil
}

func (r *PostgresInventoryRepository) GetAll(ctx context.Context, organizationID uuid.UUID) ([]entity.Ingredient, error) {
	query := `SELECT id, organization_id, name, unit, quantity, minimum_quantity, created_at, updated_at FROM ingredients WHERE organization_id = $1`

	rows, err := r.db.Query(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []entity.Ingredient
	for rows.Next() {
		var ingredient entity.Ingredient
		err = rows.Scan(&ingredient.ID, &ingredient.OrganizationID, &ingredient.Name, &ingredient.Unit, &ingredient.Quantity, &ingredient.MinimumQuantity, &ingredient.CreatedAt, &ingredient.UpdatedAt)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (r *PostgresInventoryRepository) Create(ctx context.Context, ingredient entity.Ingredient) error {
	query := `INSERT INTO ingredients (id, organization_id, name, unit, quantity, minimum_quantity, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query, ingredient.ID, ingredient.OrganizationID, ingredient.Name, ingredient.Unit, ingredient.Quantity, ingredient.MinimumQuantity, ingredient.CreatedAt, ingredient.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresInventoryRepository) CreateMany(ctx context.Context, ingredients []entity.Ingredient) error {
	batch := &pgx.Batch{}
	for _, ingredient := range ingredients {
		batch.Queue(`INSERT INTO ingredients (organization_id, name, unit, quantity, minimum_quantity, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`, ingredient.OrganizationID, ingredient.Name, ingredient.Unit, ingredient.Quantity, ingredient.MinimumQuantity, ingredient.CreatedAt, ingredient.UpdatedAt)
	}

	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("batch insert failed for query %d: %w", i, err)
		}
	}

	return br.Close()
}

func (r *PostgresInventoryRepository) Delete(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) error {
	query := `DELETE FROM ingredients WHERE id = $1 AND organization_id = $2`
	_, err := r.db.Exec(ctx, query, id, organizationID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresInventoryRepository) DeleteMany(ctx context.Context, ids []uuid.UUID, organizationID uuid.UUID) error {
	batch := &pgx.Batch{}
	for _, id := range ids {
		batch.Queue(`DELETE FROM ingredients WHERE id = $1 AND organization_id = $2`, id, organizationID)
	}

	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("batch delete failed for query %d: %w", i, err)
		}
	}

	return br.Close()
}

func (r *PostgresInventoryRepository) Update(ctx context.Context, ingredient entity.Ingredient) error {
	query := `UPDATE ingredients SET`
	args := []interface{}{}
	argCount := 1

	if ingredient.Name != "" {
		query += fmt.Sprintf(" name = $%v", argCount)
		args = append(args, ingredient.Name)
		argCount++
	}
	if ingredient.Quantity > 0 {
		query += fmt.Sprintf(" quantity = $%v", argCount)
		args = append(args, ingredient.Quantity)
		argCount++
	}
	if ingredient.MinimumQuantity > 0 {
		query += fmt.Sprintf(" minimum_quantity = $%v", argCount)
		args = append(args, ingredient.MinimumQuantity)
		argCount++
	}
	if ingredient.Unit != "" {
		query += fmt.Sprintf(" unit = $%v", argCount)
		args = append(args, ingredient.Unit)
		argCount++
	}

	query += fmt.Sprintf(" updated_at = $%v", argCount)
	args = append(args, time.Now())
	argCount++

	query += fmt.Sprintf(" WHERE id = $%v", argCount)
	args = append(args, ingredient.ID)

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

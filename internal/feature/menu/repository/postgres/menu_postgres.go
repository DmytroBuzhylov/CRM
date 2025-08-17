package postgres

import (
	"Test/internal/feature/menu/entity"
	"Test/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"time"
)

type PostgresMenuRepository struct {
	db *pgxpool.Pool
}

func NewPostgresMenuRepository(db *pgxpool.Pool) *PostgresMenuRepository {
	return &PostgresMenuRepository{db: db}
}

type PostgresMenuTransactionalRepository struct {
	pool *pgxpool.Pool
	Tx   pgx.Tx
}

func NewPostgresMenuTransactionalRepository(pool *pgxpool.Pool) *PostgresMenuTransactionalRepository {
	if pool == nil {
		return nil
	}
	return &PostgresMenuTransactionalRepository{pool: pool}
}

func (r *PostgresMenuTransactionalRepository) CreateMenuItem(ctx context.Context, menuItem entity.MenuItem) error {
	query := `INSERT INTO menu_items (id, organization_id, name, description, price, category, is_available, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.Tx.Exec(ctx, query,
		menuItem.ID,
		menuItem.OrganizationID,
		menuItem.Name,
		menuItem.Description,
		menuItem.Price.Decimal,
		menuItem.Category,
		menuItem.IsAvailable,
		menuItem.UpdatedAt,
		menuItem.CreatedAt,
	)

	return err
}

// WithTx Must be used before calling PostgresMenuTransactionalRepository methods
func (r *PostgresMenuTransactionalRepository) WithTx(ctx context.Context) (*PostgresMenuTransactionalRepository, error) {
	//tx, ok := ctx.Value("Tx").(*pgxpool.Tx)
	//if !ok || tx == nil {
	//	return nil, errors.New("error getting transaction dependency")
	//}
	//
	//return &PostgresMenuTransactionalRepository{tx: tx}, nil

	if r.pool == nil {
		return nil, errors.New("pool is not initialized")
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return &PostgresMenuTransactionalRepository{pool: r.pool, Tx: tx}, nil
}

func (r *PostgresMenuRepository) GetMenuItem(ctx context.Context, menuItemID uuid.UUID) (item entity.MenuItem, err error) {
	query := `SELECT (id, organization_id, name, description, price, category, is_available, updated_at, created_at) FROM menu_items WHERE id = $1`

	err = r.db.QueryRow(ctx, query, menuItemID).Scan(
		&item.ID,
		&item.OrganizationID,
		&item.Name,
		&item.Description,
		&item.Price.Decimal,
		&item.Category,
		&item.IsAvailable,
		&item.UpdatedAt,
		&item.CreatedAt,
	)

	if err != nil {
		return entity.MenuItem{}, err
	}

	return item, nil
}

func (r *PostgresMenuRepository) GetAllMenuItems(ctx context.Context, organizationID uuid.UUID) (items []entity.MenuItem, err error) {
	query := `SELECT (id, organization_id, name, description, price, category, is_available, updated_at, created_at) FROM menu_items WHERE organization_id = $1`

	rows, err := r.db.Query(ctx, query, organizationID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.MenuItem
		err = rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.Name,
			&item.Description,
			&item.Price.Decimal,
			&item.Category,
			&item.IsAvailable,
			&item.UpdatedAt,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return items, nil
}

func (r *PostgresMenuTransactionalRepository) UpdateMenuItem(ctx context.Context, menuItem entity.MenuItem) error {
	query := `UPDATE menu_items SET`
	var (
		args      = []interface{}{}
		count int = 1
	)

	if menuItem.Name != "" {
		if len(args) > 0 {
			query += " AND"
		}
		query += fmt.Sprintf(" name = $%v", count)
		args = append(args, menuItem.Name)
		count++
	}
	if menuItem.Price.Valid {
		if len(args) > 0 {
			query += " AND"
		}
		query += fmt.Sprintf(" price = $%v", count)
		args = append(args, menuItem.Price.Decimal)
		count++
	}
	if menuItem.Category != "" {
		if len(args) > 0 {
			query += " AND"
		}
		query += fmt.Sprintf(" category = $%v", count)
		args = append(args, menuItem.Category)
		count++
	}
	if menuItem.Description != "" {
		if len(args) > 0 {
			query += " AND"
		}
		query += fmt.Sprintf(" description = $%v", count)
		args = append(args, menuItem.Description)
		count++
	}
	if len(args) > 0 {
		query += " AND"
	}
	query += fmt.Sprintf(" updated_at = $%v", count)
	args = append(args, time.Now())
	count++

	query += fmt.Sprintf(" WHERE id = $%v", count)
	args = append(args, menuItem.ID)

	_, err := r.Tx.Exec(ctx, query, args...)

	return err
}

func (r *PostgresMenuTransactionalRepository) DeleteMenuItem(ctx context.Context, menuItemID uuid.UUID, organizationID uuid.UUID) error {
	query := `DELETE FROM menu_items WHERE id = $1 AND organization_id = $2`
	_, err := r.Tx.Exec(ctx, query, menuItemID, organizationID)

	return err
}

func (r *PostgresMenuTransactionalRepository) AddRecipeItem(ctx context.Context, recipeItem entity.RecipeItem) error {
	query := `INSERT INTO recipe_items (id, menu_item_id, ingredient_id, quantity_needed, unit_of_measure, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.Tx.Exec(ctx, query,
		&recipeItem.ID,
		&recipeItem.MenuItemID,
		&recipeItem.IngredientID,
		&recipeItem.QuantityNeeded,
		&recipeItem.UnitOfMeasure,
		&recipeItem.UpdatedAt,
		&recipeItem.CreatedAt,
	)

	return err
}

func (r *PostgresMenuTransactionalRepository) AddRecipeItems(ctx context.Context, recipeItems []entity.RecipeItem) error {
	batch := &pgx.Batch{}
	for _, i := range recipeItems {
		batch.Queue(`INSERT INTO recipe_items (menu_item_id, ingredient_id, quantity_needed, unit_of_measure, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
			&i.MenuItemID,
			&i.IngredientID,
			&i.QuantityNeeded,
			&i.UnitOfMeasure,
			&i.UpdatedAt,
			&i.CreatedAt,
		)
	}

	br := r.Tx.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("batch delete failed for query %d: %w", i, err)
		}
	}

	return br.Close()
}

func (r *PostgresMenuTransactionalRepository) UpdateRecipeItem(ctx context.Context, recipeItem entity.RecipeItem) error {
	query := `UPDATE recipe_items SET`
	args := []interface{}{}
	count := 1

	if recipeItem.UnitOfMeasure != "" {
		if len(args) > 0 {
			query += " AND"
		}
		if err := utils.ValidateProductUnit(recipeItem.UnitOfMeasure); err == nil {
			query += fmt.Sprintf(" unit_of_measure = $%v", count)
			args = append(args, recipeItem.UnitOfMeasure)
			count++
		}
	}
	if recipeItem.QuantityNeeded > 0 {
		if len(args) > 0 {
			query += " AND"
		}
		query += fmt.Sprintf(" quantity_needed = $%v", count)
		args = append(args, recipeItem.QuantityNeeded)
		count++
	}
	if len(args) > 0 {
		query += " AND"
	}
	query += fmt.Sprintf(" updated_at = $%v", count)
	args = append(args, time.Now())
	count++
	query += fmt.Sprintf(" WHERE id = $%v", count)
	args = append(args, recipeItem.ID)

	_, err := r.Tx.Exec(ctx, query, args...)

	return err
}

func (r *PostgresMenuTransactionalRepository) DeleteRecipeItem(ctx context.Context, itemID uuid.UUID, organizationID uuid.UUID) error {
	query := `DELETE FROM recipe_items WHERE id = $1`

	_, err := r.Tx.Exec(ctx, query, itemID)

	return err
}

func (r *PostgresMenuTransactionalRepository) DecreaseInventory(ctx context.Context, quantity decimal.NullDecimal, ingredientID uuid.UUID, organizationID uuid.UUID) error {
	query := `UPDATE ingredients SET quantity = quantity - $1 WHERE id = $2 AND organization_id = $3`

	_, err := r.Tx.Exec(ctx, query, quantity.Decimal, ingredientID, organizationID)

	return err
}

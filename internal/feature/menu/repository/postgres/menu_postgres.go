package postgres

import (
	entity2 "Test/internal/feature/inventory/entity"
	"Test/internal/feature/menu/entity"
	"Test/internal/feature/menu/repository"
	"Test/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"time"
)

type PostgresMenuRepository struct {
	db *pgxpool.Pool
}

func NewPostgresMenuRepository(db *pgxpool.Pool) *PostgresMenuRepository {
	return &PostgresMenuRepository{db: db}
}

type PostgresMenuTransactionalRepository struct {
	tx pgx.Tx
}

func (r *PostgresMenuTransactionalRepository) CreateMenuItem(ctx context.Context, menuItem entity.MenuItem) error {
	query := `INSERT INTO menu_items (id, organization_id, name, description, price, category, is_available, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.tx.Exec(ctx, query,
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

func (r *PostgresMenuRepository) StartTx(ctx context.Context, opts pgx.TxOptions) (repository.MenuTransactionalRepository, error) {
	if r.db == nil {
		return nil, errors.New("pgxpool nil")
	}
	tx, err := r.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &PostgresMenuTransactionalRepository{tx: tx}, nil
}

func (r *PostgresMenuTransactionalRepository) Rollback(ctx context.Context) error {
	return r.tx.Rollback(ctx)
}

func (r *PostgresMenuTransactionalRepository) Commit(ctx context.Context) error {
	return r.tx.Commit(ctx)
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

	_, err := r.tx.Exec(ctx, query, args...)

	return err
}

func (r *PostgresMenuTransactionalRepository) DeleteMenuItem(ctx context.Context, menuItemID uuid.UUID, organizationID uuid.UUID) error {
	query := `DELETE FROM menu_items WHERE id = $1 AND organization_id = $2`
	_, err := r.tx.Exec(ctx, query, menuItemID, organizationID)

	return err
}

func (r *PostgresMenuTransactionalRepository) AddRecipeItem(ctx context.Context, recipeItem entity.RecipeItem) error {
	query := `INSERT INTO recipe_items (id, menu_item_id, ingredient_id, quantity_needed, unit_of_measure, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.tx.Exec(ctx, query,
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

	br := r.tx.SendBatch(ctx, batch)
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

	_, err := r.tx.Exec(ctx, query, args...)

	return err
}

func (r *PostgresMenuTransactionalRepository) DeleteRecipeItem(ctx context.Context, itemID uuid.UUID, organizationID uuid.UUID) error {
	query := `DELETE FROM recipe_items WHERE id = $1`

	_, err := r.tx.Exec(ctx, query, itemID)

	return err
}

func (r *PostgresMenuTransactionalRepository) DecreaseInventory(ctx context.Context, menuItemID uuid.UUID, quantity uint, organizationID uuid.UUID) error {
	query := `SELECT EXISTS (SELECT 1 FROM menu_items WHERE id = $1 AND organization_id = $2)`
	var exist bool
	err := r.tx.QueryRow(ctx, query, menuItemID, organizationID).Scan(&exist)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("the dish does not belong to this organization")
	}

	query = `SELECT id, ingredient_id, quantity_needed, unit_of_measure FROM recipe_items WHERE menu_item_id = $1`
	rows, err := r.tx.Query(ctx, query, menuItemID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var (
		recipeItems    []entity.RecipeItem
		recipeItemsIDs []uuid.UUID
	)

	for rows.Next() {
		var recipeItem entity.RecipeItem
		err := rows.Scan(&recipeItem.ID, &recipeItem.IngredientID, &recipeItem.QuantityNeeded, &recipeItem.UnitOfMeasure)
		if err != nil {
			return err
		}
		recipeItems = append(recipeItems, recipeItem)
		recipeItemsIDs = append(recipeItemsIDs, recipeItem.IngredientID)
	}
	if rows.Err() != nil {
		return err
	}

	query = `SELECT id, unit, quantity FROM ingredients WHERE id = ANY($1) AND organization_id = $2`
	rows, err = r.tx.Query(ctx, query, recipeItemsIDs)
	if err != nil {
		return err
	}

	defer rows.Close()

	var ingredients []entity2.Ingredient
	for rows.Next() {
		var ingredient entity2.Ingredient
		err = rows.Scan(&ingredient.ID, &ingredient.Unit, &ingredient.Quantity)
		if err != nil {
			return err
		}
		ingredients = append(ingredients, ingredient)
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	ingredientMap := make(map[uuid.UUID]entity2.Ingredient)
	for _, ing := range ingredients {
		ingredientMap[ing.ID] = ing
	}
	for _, recipeItem := range recipeItems {
		ing, ok := ingredientMap[recipeItem.IngredientID]
		if !ok {
			return errors.New("missing ingredient to recipe: " + recipeItem.IngredientID.String())
		}

		if ing.Unit != recipeItem.UnitOfMeasure {
			return errors.New("unit mismatch for ingredient " + ing.ID.String())
		}
		needed := recipeItem.QuantityNeeded * uint64(quantity)
		if ing.Quantity < needed {
			return errors.New("insufficient inventory for ingredient " + ing.ID.String() + ": need " + strconv.FormatUint(needed, 10) + ", have " + strconv.FormatUint(ing.Quantity, 10))
		}
	}

	query = `UPDATE ingredients SET quantity = quantity - $1 WHERE id = $2 AND organization_id = $3`
	batch := &pgx.Batch{}

	for _, i := range recipeItems {
		batch.Queue(query, i.QuantityNeeded*uint64(quantity), i.IngredientID, organizationID)
	}

	br := r.tx.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err = br.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

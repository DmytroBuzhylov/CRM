package postgres

import (
	"Test/internal/feature/order/entity"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type PostgresOrderRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresOrderRepository(pool *pgxpool.Pool) *PostgresOrderRepository {
	return &PostgresOrderRepository{pool: pool}
}

func (r *PostgresOrderRepository) CreateOrder(ctx context.Context, order *entity.Order, items []*entity.OrderItem) error {

	if len(items) == 0 {
		return fmt.Errorf("order must contain at least one item")
	}

	order.ID = uuid.New()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	uniqueMenuItemsMap := make(map[uuid.UUID]struct{})
	for _, i := range items {
		uniqueMenuItemsMap[i.MenuItemID] = struct{}{}
	}

	var uniqueMenuItemIDS []uuid.UUID
	for id := range uniqueMenuItemsMap {
		uniqueMenuItemIDS = append(uniqueMenuItemIDS, id)
	}

	query := `SELECT id, price FROM menu_items WHERE id = ANY($1) AND organization_id = $2`
	rows, err := tx.Query(ctx, query, uniqueMenuItemIDS, order.OrganizationID)
	if err != nil {
		return err
	}
	defer rows.Close()

	menuItemsPrice := make(map[uuid.UUID]decimal.Decimal)
	for rows.Next() {
		var (
			id    uuid.UUID
			price decimal.Decimal
		)
		err = rows.Scan(&id, &price)
		if err != nil {
			return err
		}
		menuItemsPrice[id] = price
	}
	fmt.Println(menuItemsPrice)
	if rows.Err() != nil {
		return rows.Err()
	}
	if len(menuItemsPrice) != len(uniqueMenuItemIDS) {
		return fmt.Errorf("not all menu items were found for the provided IDs")
	}

	var totalPrice decimal.Decimal
	for _, i := range items {
		price, exist := menuItemsPrice[i.MenuItemID]
		if !exist {
			return fmt.Errorf("menu item %s not found in price map", i.MenuItemID)
		}
		totalPrice = totalPrice.Add(price.Mul(decimal.NewFromUint64(i.Quantity)))
	}

	query = `INSERT INTO orders (id, organization_id, user_id, table_id, status, user_type, notes, total_price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.Exec(ctx, query,
		order.ID,
		order.OrganizationID,
		order.UserID,
		order.TableID,
		"pending",
		order.UserType,
		order.Notes,
		totalPrice,
	)
	if err != nil {
		return err
	}

	query = `INSERT INTO order_items (order_id, menu_item_id, quantity, price_at_time_of_order) VALUES ($1, $2, $3, $4)`

	batch := &pgx.Batch{}
	for _, i := range items {
		batch.Queue(query, order.ID, i.MenuItemID, i.Quantity, menuItemsPrice[i.MenuItemID])
	}

	br := tx.SendBatch(ctx, batch)

	for i := 0; i < batch.Len(); i++ {
		_, err = br.Exec()
		if err != nil {
			br.Close()
			return fmt.Errorf("failed to execute batch operation %d: %w", i, err)
		}
	}
	br.Close()

	return tx.Commit(ctx)
}

func (r *PostgresOrderRepository) UpdateOrderStatus(ctx context.Context, id uuid.UUID, organizationID uuid.UUID, status string) error {

	query := `UPDATE orders SET status = $3 WHERE id = $1 AND organization_id = $2`
	_, err := r.pool.Exec(ctx, query, id, organizationID, status)

	return err
}

func (r *PostgresOrderRepository) UpdateOrderNotes(ctx context.Context, id uuid.UUID, organizationID uuid.UUID, notes string) error {

	query := `UPDATE orders SET notes = $3 WHERE id = $1 AND organization_id = $2`
	_, err := r.pool.Exec(ctx, query, id, organizationID, notes)

	return err
}

func (r *PostgresOrderRepository) GetOrder(ctx context.Context, id uuid.UUID) (*entity.Order, error) {

	query := `SELECT organization_id, user_id, table_id, user_type, status, notes, total_price, created_at, updated_at FROM orders WHERE id = $1`
	var order *entity.Order
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&order.OrganizationID,
		&order.UserID,
		&order.TableID,
		&order.UserType,
		&order.Status,
		&order.Notes,
		&order.TotalPrice,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	order.ID = id

	return order, nil
}

func (r *PostgresOrderRepository) GetAllOrders(ctx context.Context, organizationID uuid.UUID) ([]*entity.Order, error) {

	query := `SELECT id, user_id, table_id, user_type, status, notes, total_price, created_at, updated_at FROM orders WHERE organization_id = $1`
	rows, err := r.pool.Query(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]*entity.Order, 0)
	for rows.Next() {

		var order *entity.Order
		order.OrganizationID = organizationID

		err = rows.Scan(
			&order.ID,
			&order.UserID,
			&order.TableID,
			&order.UserType,
			&order.Status,
			&order.Notes,
			&order.TotalPrice,
			&order.CreatedAt,
			&order.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *PostgresOrderRepository) GetOrderWithItems(ctx context.Context, id uuid.UUID) (*entity.OrderWithItems, error) {

	query := `
        SELECT id, organization_id, user_id, table_id, user_type, status, notes, total_price, created_at, updated_at
        FROM orders
        WHERE id = $1
    `
	var order entity.Order
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.OrganizationID,
		&order.UserID,
		&order.TableID,
		&order.UserType,
		&order.Status,
		&order.Notes,
		&order.TotalPrice,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {

		return nil, err
	}

	query = `
        SELECT id, order_id, menu_item_id, quantity, price_at_time_of_order, created_at, updated_at
        FROM order_items
        WHERE order_id = $1
    `
	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*entity.OrderItem
	for rows.Next() {
		var item entity.OrderItem
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.MenuItemID,
			&item.Quantity,
			&item.PriceAtTimeOfOrder,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &entity.OrderWithItems{
		Order: &order,
		Items: items,
	}, nil
}

func (r *PostgresOrderRepository) GetAllOrdersWithItems(ctx context.Context, organizationID uuid.UUID) ([]*entity.OrderWithItems, error) {

	query := `
	SELECT id, user_id, table_id, user_type, status, notes, total_price, created_at, updated_at
        FROM orders
        WHERE organization_id = $1
	`

	rows, err := r.pool.Query(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		orders   = make([]*entity.Order, 0)
		orderIDs []uuid.UUID
	)

	for rows.Next() {
		var order = new(entity.Order)
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.TableID,
			&order.UserType,
			&order.Status,
			&order.Notes,
			&order.TotalPrice,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		order.OrganizationID = organizationID
		orders = append(orders, order)
		orderIDs = append(orderIDs, order.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	orderItemsMap := make(map[uuid.UUID][]*entity.OrderItem, len(orders))

	query = `
		SELECT id, order_id, menu_item_id, quantity, price_at_time_of_order, created_at, updated_at
        FROM order_items
        WHERE order_id = ANY($1)
	`

	rows, err = r.pool.Query(ctx, query, orderIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item = new(entity.OrderItem)
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.MenuItemID,
			&item.Quantity,
			&item.PriceAtTimeOfOrder,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orderItemsMap[item.OrderID] = append(orderItemsMap[item.OrderID], item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var respDTO = make([]*entity.OrderWithItems, 0, len(orderIDs))
	for _, o := range orders {
		respDTO = append(respDTO, &entity.OrderWithItems{
			Order: o,
			Items: orderItemsMap[o.ID],
		})
	}
	return respDTO, nil
}

func (r *PostgresOrderRepository) GetOrderItem(ctx context.Context, orderId uuid.UUID) (*entity.OrderItem, error) {

	query := `
		SELECT id, order_id, menu_item_id, quantity, price_at_time_of_order, created_at, updated_at
        FROM order_items
        WHERE id = $1
	`
	var item *entity.OrderItem
	err := r.pool.QueryRow(ctx, query, orderId).Scan(
		&item.ID,
		&item.OrderID,
		&item.MenuItemID,
		&item.Quantity,
		&item.PriceAtTimeOfOrder,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *PostgresOrderRepository) GetOrderItems(ctx context.Context, ids []uuid.UUID) ([]*entity.OrderItem, error) {
	query := `
		SELECT id, order_id, menu_item_id, quantity, price_at_time_of_order, created_at, updated_at
        FROM order_items
        WHERE id = $1
	`
	rows, err := r.pool.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items = make([]*entity.OrderItem, 0, len(ids))
	for rows.Next() {
		var item *entity.OrderItem

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.MenuItemID,
			&item.Quantity,
			&item.PriceAtTimeOfOrder,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *PostgresOrderRepository) CreateOrderItems(ctx context.Context, organizationID uuid.UUID, orderID uuid.UUID, items []*entity.OrderItem) error {
	query := `
	SELECT EXISTS(
    SELECT 1
    FROM orders
    WHERE organization_id = $1
      AND id = $2
	)
	`
	var exists bool
	err := r.pool.QueryRow(ctx, query, organizationID, organizationID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("order not found")
	}

	query = `INSERT INTO order_items VALUES (id, order_id, menu_item_id, quantity, price_at_time_of_order, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7)`

	batch := &pgx.Batch{}
	for _, item := range items {
		batch.Queue(query,
			item.ID,
			item.OrderID,
			item.MenuItemID,
			item.Quantity,
			item.PriceAtTimeOfOrder,
			item.CreatedAt,
			item.UpdatedAt,
		)
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()

	for range items {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return nil
}

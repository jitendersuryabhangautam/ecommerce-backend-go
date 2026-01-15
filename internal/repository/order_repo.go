package repository

import (
	"context"
	"errors"
	"fmt"

	"ecommerce-backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	CreateWithTx(ctx context.Context, tx pgx.Tx, order *models.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	GetByOrderNumber(ctx context.Context, orderNumber string) (*models.Order, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Order, int, error)
	GetAll(ctx context.Context, page, limit int, status string) ([]models.Order, int, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.OrderStatus) error
	CancelOrder(ctx context.Context, id uuid.UUID) error
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	err = r.CreateWithTx(ctx, tx, order)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *orderRepository) CreateWithTx(ctx context.Context, tx pgx.Tx, order *models.Order) error {
	// Insert order
	orderQuery := `
        INSERT INTO orders (id, user_id, order_number, total_amount, status, 
                          shipping_address, billing_address)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING created_at, updated_at
    `

	_, err := tx.Exec(ctx, orderQuery,
		order.ID,
		order.UserID,
		order.OrderNumber,
		order.TotalAmount,
		order.Status,
		order.ShippingAddress,
		order.BillingAddress,
	)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// Insert order items
	itemQuery := `
        INSERT INTO order_items (id, order_id, product_id, quantity, price_at_time)
        VALUES ($1, $2, $3, $4, $5)
    `

	for _, item := range order.Items {
		_, err := tx.Exec(ctx, itemQuery,
			uuid.New(),
			order.ID,
			item.ProductID,
			item.Quantity,
			item.PriceAtTime,
		)

		if err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}
	}

	return nil
}

func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	// Get order
	orderQuery := `
        SELECT id, user_id, order_number, total_amount, status, 
               shipping_address, billing_address, created_at, updated_at
        FROM orders
        WHERE id = $1
    `

	var order models.Order
	err := r.db.QueryRow(ctx, orderQuery, id).Scan(
		&order.ID,
		&order.UserID,
		&order.OrderNumber,
		&order.TotalAmount,
		&order.Status,
		&order.ShippingAddress,
		&order.BillingAddress,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	// Get order items
	itemsQuery := `
        SELECT 
            oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price_at_time, oi.created_at,
            p.id, p.sku, p.name, p.description, p.price, p.stock_quantity, 
            p.category, p.image_url, p.created_at, p.updated_at
        FROM order_items oi
        JOIN products p ON oi.product_id = p.id
        WHERE oi.order_id = $1
        ORDER BY oi.created_at
    `

	rows, err := r.db.Query(ctx, itemsQuery, order.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		var product models.Product

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.PriceAtTime,
			&item.CreatedAt,
			&product.ID,
			&product.SKU,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.Category,
			&product.ImageURL,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		item.Product = product
		items = append(items, item)
	}

	order.Items = items
	return &order, nil
}

func (r *orderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*models.Order, error) {
	orderQuery := `
        SELECT id, user_id, order_number, total_amount, status, 
               shipping_address, billing_address, created_at, updated_at
        FROM orders
        WHERE order_number = $1
    `

	var order models.Order
	err := r.db.QueryRow(ctx, orderQuery, orderNumber).Scan(
		&order.ID,
		&order.UserID,
		&order.OrderNumber,
		&order.TotalAmount,
		&order.Status,
		&order.ShippingAddress,
		&order.BillingAddress,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Order, int, error) {
	offset := (page - 1) * limit

	// Count total orders
	countQuery := `SELECT COUNT(*) FROM orders WHERE user_id = $1`
	var total int
	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get orders with pagination
	ordersQuery := `
        SELECT id, user_id, order_number, total_amount, status, 
               shipping_address, billing_address, created_at, updated_at
        FROM orders
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.db.Query(ctx, ordersQuery, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderNumber,
			&order.TotalAmount,
			&order.Status,
			&order.ShippingAddress,
			&order.BillingAddress,
			&order.CreatedAt,
			&order.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		orders = append(orders, order)
	}

	return orders, total, nil
}

func (r *orderRepository) GetAll(ctx context.Context, page, limit int, status string) ([]models.Order, int, error) {
	offset := (page - 1) * limit

	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	// Count total orders
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM orders %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get orders with pagination
	ordersQuery := fmt.Sprintf(`
        SELECT id, user_id, order_number, total_amount, status, 
               shipping_address, billing_address, created_at, updated_at
        FROM orders %s
        ORDER BY created_at DESC
        LIMIT $%d OFFSET $%d
    `, whereClause, argCount, argCount+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, ordersQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderNumber,
			&order.TotalAmount,
			&order.Status,
			&order.ShippingAddress,
			&order.BillingAddress,
			&order.CreatedAt,
			&order.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		orders = append(orders, order)
	}

	return orders, total, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.OrderStatus) error {
	query := `
        UPDATE orders
        SET status = $1, updated_at = NOW()
        WHERE id = $2
    `

	result, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("order not found")
	}

	return nil
}

func (r *orderRepository) CancelOrder(ctx context.Context, id uuid.UUID) error {
	query := `
        UPDATE orders
        SET status = 'cancelled', updated_at = NOW()
        WHERE id = $1 AND status IN ('pending', 'processing')
    `

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("order cannot be cancelled or not found")
	}

	return nil
}

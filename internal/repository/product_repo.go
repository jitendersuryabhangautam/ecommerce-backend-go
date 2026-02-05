package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"ecommerce-backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetBySKU(ctx context.Context, sku string) (*models.Product, error)
	GetAll(ctx context.Context, page, limit int, category, search string) ([]models.Product, int, error)
	GetAllAdmin(ctx context.Context, page, limit, rangeDays int) ([]models.Product, int, error)
	GetTopProducts(ctx context.Context, limit, rangeDays int) ([]models.TopProductItem, error)
	Update(ctx context.Context, id uuid.UUID, updateData *models.ProductUpdateRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error
	UpdateStockWithTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, quantity int) error
	GetStock(ctx context.Context, id uuid.UUID) (int, error)
	ReserveStock(ctx context.Context, productID, cartID uuid.UUID, quantity int, expiresAt int64) error
	ReleaseStockReservation(ctx context.Context, productID, cartID uuid.UUID) error
	GetAvailableStock(ctx context.Context, productID uuid.UUID) (int, error)
	GetAvailableStockExcludingCart(ctx context.Context, productID, cartID uuid.UUID) (int, error)
}

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	query := `
        INSERT INTO products (sku, name, description, price, stock_quantity, category, image_url)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
    `

	return r.db.QueryRow(ctx, query,
		product.SKU,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.Category,
		product.ImageURL,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	query := `
        SELECT 
            p.id, p.sku, p.name, p.description, p.price,
            p.stock_quantity - COALESCE(SUM(sr.quantity), 0) as available_stock,
            p.category, p.image_url, p.created_at, p.updated_at
        FROM products p
        LEFT JOIN stock_reservations sr ON p.id = sr.product_id 
            AND sr.expires_at > NOW()
        WHERE p.id = $1
        GROUP BY p.id, p.sku, p.name, p.description, p.price, p.stock_quantity,
                 p.category, p.image_url, p.created_at, p.updated_at
    `

	var product models.Product
	err := r.db.QueryRow(ctx, query, id).Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*models.Product, error) {
	query := `
        SELECT 
            p.id, p.sku, p.name, p.description, p.price,
            p.stock_quantity - COALESCE(SUM(sr.quantity), 0) as available_stock,
            p.category, p.image_url, p.created_at, p.updated_at
        FROM products p
        LEFT JOIN stock_reservations sr ON p.id = sr.product_id 
            AND sr.expires_at > NOW()
        WHERE p.sku = $1
        GROUP BY p.id, p.sku, p.name, p.description, p.price, p.stock_quantity,
                 p.category, p.image_url, p.created_at, p.updated_at
    `

	var product models.Product
	err := r.db.QueryRow(ctx, query, sku).Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) GetAll(ctx context.Context, page, limit int, category, search string) ([]models.Product, int, error) {
	offset := (page - 1) * limit

	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if category != "" {
		whereClause += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, category)
		argCount++
	}

	if search != "" {
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+search+"%")
		argCount++
	}

	// Count total products
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get products with pagination
	productsQuery := fmt.Sprintf(`
        SELECT 
            p.id, p.sku, p.name, p.description, p.price, 
            p.stock_quantity - COALESCE(SUM(sr.quantity), 0) as available_stock,
            p.category, p.image_url, p.created_at, p.updated_at
        FROM products p
        LEFT JOIN stock_reservations sr ON p.id = sr.product_id 
            AND sr.expires_at > NOW()
        %s
        GROUP BY p.id, p.sku, p.name, p.description, p.price, p.stock_quantity, 
                 p.category, p.image_url, p.created_at, p.updated_at
        ORDER BY p.created_at DESC
        LIMIT $%d OFFSET $%d
    `, whereClause, argCount, argCount+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, productsQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
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
			return nil, 0, err
		}
		products = append(products, product)
	}

	return products, total, nil
}

func (r *productRepository) GetAllAdmin(ctx context.Context, page, limit, rangeDays int) ([]models.Product, int, error) {
	offset := (page - 1) * limit

	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if rangeDays > 0 {
		whereClause += fmt.Sprintf(" AND created_at >= NOW() - $%d * INTERVAL '1 day'", argCount)
		args = append(args, rangeDays)
		argCount++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
        SELECT 
            p.id, p.sku, p.name, p.description, p.price,
            p.stock_quantity - COALESCE(SUM(sr.quantity), 0) as available_stock,
            p.category, p.image_url, p.created_at, p.updated_at
        FROM products p
        LEFT JOIN stock_reservations sr ON p.id = sr.product_id 
            AND sr.expires_at > NOW()
        %s
        GROUP BY p.id, p.sku, p.name, p.description, p.price, p.stock_quantity,
                 p.category, p.image_url, p.created_at, p.updated_at
        ORDER BY p.created_at DESC
        LIMIT $%d OFFSET $%d
    `, whereClause, argCount, argCount+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
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
			return nil, 0, err
		}
		products = append(products, product)
	}

	return products, total, nil
}

func (r *productRepository) GetTopProducts(ctx context.Context, limit, rangeDays int) ([]models.TopProductItem, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if rangeDays > 0 {
		whereClause += fmt.Sprintf(" AND o.created_at >= NOW() - $%d * INTERVAL '1 day'", argCount)
		args = append(args, rangeDays)
		argCount++
	}

	query := fmt.Sprintf(`
        SELECT 
            p.id, p.sku, p.name, p.description, p.price, p.stock_quantity, p.category, p.image_url,
            p.created_at, p.updated_at,
            COALESCE(SUM(oi.quantity), 0) as total_quantity,
            COALESCE(SUM(oi.quantity * oi.price_at_time), 0) as total_revenue
        FROM order_items oi
        JOIN orders o ON oi.order_id = o.id
        JOIN products p ON oi.product_id = p.id
        %s
        GROUP BY p.id
        ORDER BY total_quantity DESC, total_revenue DESC
        LIMIT $%d
    `, whereClause, argCount)

	args = append(args, limit)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.TopProductItem
	for rows.Next() {
		var item models.TopProductItem
		var product models.Product
		if err := rows.Scan(
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
			&item.TotalQuantity,
			&item.TotalRevenue,
		); err != nil {
			return nil, err
		}
		item.Product = product
		items = append(items, item)
	}

	return items, nil
}

func (r *productRepository) Update(ctx context.Context, id uuid.UUID, updateData *models.ProductUpdateRequest) error {
	query := "UPDATE products SET "
	args := []interface{}{}
	argCount := 1

	updates := []string{}

	if updateData.Name != "" {
		updates = append(updates, fmt.Sprintf("name = $%d", argCount))
		args = append(args, updateData.Name)
		argCount++
	}

	if updateData.Description != "" {
		updates = append(updates, fmt.Sprintf("description = $%d", argCount))
		args = append(args, updateData.Description)
		argCount++
	}

	if updateData.Price > 0 {
		updates = append(updates, fmt.Sprintf("price = $%d", argCount))
		args = append(args, updateData.Price)
		argCount++
	}

	if updateData.Stock >= 0 {
		updates = append(updates, fmt.Sprintf("stock_quantity = $%d", argCount))
		args = append(args, updateData.Stock)
		argCount++
	}

	if updateData.Category != "" {
		updates = append(updates, fmt.Sprintf("category = $%d", argCount))
		args = append(args, updateData.Category)
		argCount++
	}

	if updateData.ImageURL != "" {
		updates = append(updates, fmt.Sprintf("image_url = $%d", argCount))
		args = append(args, updateData.ImageURL)
		argCount++
	}

	if len(updates) == 0 {
		return nil // Nothing to update
	}

	updates = append(updates, "updated_at = NOW()")
	query += strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, id)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM products WHERE id = $1"
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *productRepository) UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error {
	query := `
        UPDATE products 
        SET stock_quantity = stock_quantity + $1, updated_at = NOW()
        WHERE id = $2 AND stock_quantity + $1 >= 0
        RETURNING stock_quantity
    `

	var newStock int
	err := r.db.QueryRow(ctx, query, quantity, id).Scan(&newStock)
	return err
}

func (r *productRepository) UpdateStockWithTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, quantity int) error {
	query := `
        UPDATE products 
        SET stock_quantity = stock_quantity + $1, updated_at = NOW()
        WHERE id = $2 AND stock_quantity + $1 >= 0
        RETURNING stock_quantity
    `

	var newStock int
	err := tx.QueryRow(ctx, query, quantity, id).Scan(&newStock)
	return err
}

func (r *productRepository) GetStock(ctx context.Context, id uuid.UUID) (int, error) {
	query := "SELECT stock_quantity FROM products WHERE id = $1"

	var stock int
	err := r.db.QueryRow(ctx, query, id).Scan(&stock)
	if err != nil {
		return 0, err
	}

	return stock, nil
}

func (r *productRepository) ReserveStock(ctx context.Context, productID, cartID uuid.UUID, quantity int, expiresAt int64) error {
	// Use PostgreSQL advisory lock to prevent race conditions
	lockQuery := `SELECT pg_advisory_xact_lock(hashtext($1))`
	lockKey := fmt.Sprintf("product_%s", productID.String())

	_, err := r.db.Exec(ctx, lockQuery, lockKey)
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}

	// First check if reservation exists and get current quantity
	checkQuery := `SELECT quantity FROM stock_reservations WHERE product_id = $1::uuid AND cart_id = $2::uuid`
	var currentQuantity int
	err = r.db.QueryRow(ctx, checkQuery, productID, cartID).Scan(&currentQuantity)

	if err == nil {
		// Reservation exists, check if new total would exceed available stock
		newTotalQuantity := currentQuantity + quantity

		// Check available stock with new total
		availQuery := `
			SELECT p.stock_quantity - COALESCE(SUM(sr.quantity), 0) + COALESCE((SELECT quantity FROM stock_reservations WHERE product_id = $1::uuid AND cart_id = $2::uuid), 0) as available
			FROM products p
			LEFT JOIN stock_reservations sr ON p.id = sr.product_id AND sr.expires_at > NOW()
			WHERE p.id = $1::uuid
			GROUP BY p.id, p.stock_quantity
		`
		var available int
		availErr := r.db.QueryRow(ctx, availQuery, productID, cartID).Scan(&available)
		if availErr != nil && !errors.Is(availErr, pgx.ErrNoRows) {
			return fmt.Errorf("failed to check available stock: %w", availErr)
		}

		if available < newTotalQuantity {
			return fmt.Errorf("insufficient stock available for reservation")
		}

		// Update existing reservation
		updateQuery := `
			UPDATE stock_reservations 
			SET quantity = $1::integer, expires_at = to_timestamp($2)
			WHERE product_id = $3::uuid AND cart_id = $4::uuid
		`
		_, err := r.db.Exec(ctx, updateQuery, newTotalQuantity, expiresAt, productID, cartID)
		return err
	} else if errors.Is(err, pgx.ErrNoRows) {
		// No existing reservation, insert new one
		// Check available stock for new reservation
		query := `
			WITH available_stock AS (
				SELECT 
					p.stock_quantity - COALESCE(SUM(sr.quantity), 0) as available
				FROM products p
				LEFT JOIN stock_reservations sr ON p.id = sr.product_id 
					AND sr.expires_at > NOW()
				WHERE p.id = $1::uuid
				GROUP BY p.id, p.stock_quantity
			)
			INSERT INTO stock_reservations (product_id, cart_id, quantity, expires_at)
			SELECT $1::uuid, $2::uuid, $3::integer, to_timestamp($4)
			FROM available_stock
			WHERE available >= $3::integer
			RETURNING id
		`

		var reservationID uuid.UUID
		err = r.db.QueryRow(ctx, query, productID, cartID, quantity, expiresAt).Scan(&reservationID)

		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("insufficient stock available for reservation")
		}

		return err
	}

	return err
}

func (r *productRepository) ReleaseStockReservation(ctx context.Context, productID, cartID uuid.UUID) error {
	query := `DELETE FROM stock_reservations WHERE product_id = $1 AND cart_id = $2`
	_, err := r.db.Exec(ctx, query, productID, cartID)
	return err
}

func (r *productRepository) GetAvailableStock(ctx context.Context, productID uuid.UUID) (int, error) {
	query := `
        SELECT 
            p.stock_quantity - COALESCE(SUM(sr.quantity), 0) as available
        FROM products p
        LEFT JOIN stock_reservations sr ON p.id = sr.product_id 
            AND sr.expires_at > NOW()
        WHERE p.id = $1
        GROUP BY p.id, p.stock_quantity
    `

	var available int
	err := r.db.QueryRow(ctx, query, productID).Scan(&available)
	if err != nil {
		return 0, err
	}

	return available, nil
}

func (r *productRepository) GetAvailableStockExcludingCart(ctx context.Context, productID, cartID uuid.UUID) (int, error) {
	query := `
        SELECT 
            p.stock_quantity - COALESCE(SUM(sr.quantity), 0) as available
        FROM products p
        LEFT JOIN stock_reservations sr ON p.id = sr.product_id 
            AND sr.expires_at > NOW()
            AND sr.cart_id != $2
        WHERE p.id = $1
        GROUP BY p.id, p.stock_quantity
    `

	var available int
	err := r.db.QueryRow(ctx, query, productID, cartID).Scan(&available)
	if err != nil {
		return 0, err
	}

	return available, nil
}

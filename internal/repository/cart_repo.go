package repository

import (
	"context"
	"errors"

	"ecommerce-backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartRepository interface {
	Create(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Cart, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
	AddItem(ctx context.Context, cartID, productID uuid.UUID, quantity int) error
	UpdateItem(ctx context.Context, cartID, itemID uuid.UUID, quantity int) error
	RemoveItem(ctx context.Context, cartID, itemID uuid.UUID) error
	ClearCart(ctx context.Context, cartID uuid.UUID) error
	GetCartWithItems(ctx context.Context, cartID uuid.UUID) (*models.Cart, error)
}

type cartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) Create(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	query := `
        INSERT INTO carts (user_id)
        VALUES ($1)
        RETURNING id, created_at, updated_at
    `

	var cart models.Cart
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&cart.ID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	cart.UserID = userID
	cart.Items = []models.CartItem{}

	return &cart, nil
}

func (r *cartRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Cart, error) {
	query := `
        SELECT id, user_id, created_at, updated_at
        FROM carts
        WHERE id = $1
    `

	var cart models.Cart
	err := r.db.QueryRow(ctx, query, id).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	query := `
        SELECT id, user_id, created_at, updated_at
        FROM carts
        WHERE user_id = $1
    `

	var cart models.Cart
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		// Create cart if it doesn't exist
		return r.Create(ctx, userID)
	}

	if err != nil {
		return nil, err
	}

	// Load cart items
	cartWithItems, err := r.GetCartWithItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	return cartWithItems, nil
}

func (r *cartRepository) AddItem(ctx context.Context, cartID, productID uuid.UUID, quantity int) error {
	// First, try to get existing item
	checkQuery := `SELECT id FROM cart_items WHERE cart_id = $1 AND product_id = $2`
	var existingID string
	err := r.db.QueryRow(ctx, checkQuery, cartID, productID).Scan(&existingID)

	if err == nil {
		// Item exists, update it
		updateQuery := `
			UPDATE cart_items 
			SET quantity = quantity + $1
			WHERE cart_id = $2 AND product_id = $3
		`
		_, err := r.db.Exec(ctx, updateQuery, quantity, cartID, productID)
		return err
	} else if errors.Is(err, pgx.ErrNoRows) {
		// Item doesn't exist, insert it
		insertQuery := `
			INSERT INTO cart_items (cart_id, product_id, quantity)
			VALUES ($1, $2, $3)
		`
		_, err := r.db.Exec(ctx, insertQuery, cartID, productID, quantity)
		return err
	}

	return err
}

func (r *cartRepository) UpdateItem(ctx context.Context, cartID, itemID uuid.UUID, quantity int) error {
	query := `
        UPDATE cart_items
        SET quantity = $1
        WHERE id = $2 AND cart_id = $3
    `

	result, err := r.db.Exec(ctx, query, quantity, itemID, cartID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("cart item not found")
	}

	return nil
}

func (r *cartRepository) RemoveItem(ctx context.Context, cartID, itemID uuid.UUID) error {
	query := `
        DELETE FROM cart_items
        WHERE id = $1 AND cart_id = $2
    `

	result, err := r.db.Exec(ctx, query, itemID, cartID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("cart item not found")
	}

	return nil
}

func (r *cartRepository) ClearCart(ctx context.Context, cartID uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1`
	_, err := r.db.Exec(ctx, query, cartID)
	return err
}

func (r *cartRepository) GetCartWithItems(ctx context.Context, cartID uuid.UUID) (*models.Cart, error) {
	// Get cart
	cartQuery := `
        SELECT id, user_id, created_at, updated_at
        FROM carts
        WHERE id = $1
    `

	var cart models.Cart
	err := r.db.QueryRow(ctx, cartQuery, cartID).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Get cart items with product details
	itemsQuery := `
        SELECT 
            ci.id, ci.cart_id, ci.product_id, ci.quantity, ci.created_at,
            p.id, p.sku, p.name, p.description, p.price, p.stock_quantity, 
            p.category, p.image_url, p.created_at, p.updated_at
        FROM cart_items ci
        JOIN products p ON ci.product_id = p.id
        WHERE ci.cart_id = $1
        ORDER BY ci.created_at DESC
    `

	rows, err := r.db.Query(ctx, itemsQuery, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		var product models.Product

		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&item.Quantity,
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

	cart.Items = items
	return &cart, nil
}

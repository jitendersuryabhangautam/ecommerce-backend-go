package repository

import (
	"context"
	"fmt"

	"ecommerce-backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReturnRepository interface {
	Create(ctx context.Context, returnReq *models.Return) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Return, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Return, int, error)
	GetAll(ctx context.Context, page, limit int, status string, rangeDays int) ([]models.AdminReturn, int, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.ReturnStatus, refundAmount float64) error
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]models.Return, error)
}

type returnRepository struct {
	db *pgxpool.Pool
}

func NewReturnRepository(db *pgxpool.Pool) ReturnRepository {
	return &returnRepository{db: db}
}

func (r *returnRepository) Create(ctx context.Context, returnReq *models.Return) error {
	query := `
        INSERT INTO returns (id, order_id, user_id, reason, status, refund_amount)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING created_at, updated_at
    `

	return r.db.QueryRow(ctx, query,
		returnReq.ID,
		returnReq.OrderID,
		returnReq.UserID,
		returnReq.Reason,
		returnReq.Status,
		returnReq.RefundAmount,
	).Scan(&returnReq.CreatedAt, &returnReq.UpdatedAt)
}

func (r *returnRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Return, error) {
	query := `
        SELECT id, order_id, user_id, reason, status, refund_amount, created_at, updated_at
        FROM returns
        WHERE id = $1
    `

	var returnReq models.Return
	err := r.db.QueryRow(ctx, query, id).Scan(
		&returnReq.ID,
		&returnReq.OrderID,
		&returnReq.UserID,
		&returnReq.Reason,
		&returnReq.Status,
		&returnReq.RefundAmount,
		&returnReq.CreatedAt,
		&returnReq.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &returnReq, nil
}

func (r *returnRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Return, int, error) {
	offset := (page - 1) * limit

	// Count total returns
	countQuery := `SELECT COUNT(*) FROM returns WHERE user_id = $1`
	var total int
	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get returns with pagination
	returnsQuery := `
        SELECT id, order_id, user_id, reason, status, refund_amount, created_at, updated_at
        FROM returns
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.db.Query(ctx, returnsQuery, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var returns []models.Return
	for rows.Next() {
		var returnReq models.Return
		err := rows.Scan(
			&returnReq.ID,
			&returnReq.OrderID,
			&returnReq.UserID,
			&returnReq.Reason,
			&returnReq.Status,
			&returnReq.RefundAmount,
			&returnReq.CreatedAt,
			&returnReq.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		returns = append(returns, returnReq)
	}

	return returns, total, nil
}

func (r *returnRepository) GetAll(ctx context.Context, page, limit int, status string, rangeDays int) ([]models.AdminReturn, int, error) {
	offset := (page - 1) * limit

	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if status != "" {
		whereClause += " AND r.status = $1"
		args = append(args, status)
		argCount++
	}
	if rangeDays > 0 {
		whereClause += fmt.Sprintf(" AND r.created_at >= NOW() - $%d * INTERVAL '1 day'", argCount)
		args = append(args, rangeDays)
		argCount++
	}

	// Count total returns
	countQuery := "SELECT COUNT(*) FROM returns r " + whereClause
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get returns with pagination
	returnsQuery := `
        SELECT 
            r.id, r.order_id, r.user_id, r.reason, r.status, r.refund_amount, r.created_at, r.updated_at,
            o.order_number, u.id, u.email
        FROM returns r
        JOIN orders o ON r.order_id = o.id
        JOIN users u ON r.user_id = u.id
    ` + whereClause + ` ORDER BY r.created_at DESC LIMIT $` + fmt.Sprintf("%d", argCount) + ` OFFSET $` + fmt.Sprintf("%d", argCount+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, returnsQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var returns []models.AdminReturn
	for rows.Next() {
		var returnReq models.AdminReturn
		err := rows.Scan(
			&returnReq.ID,
			&returnReq.OrderID,
			&returnReq.Reason,
			&returnReq.Status,
			&returnReq.RefundAmount,
			&returnReq.CreatedAt,
			&returnReq.UpdatedAt,
			&returnReq.Order.OrderNumber,
			&returnReq.User.ID,
			&returnReq.User.Email,
		)

		if err != nil {
			return nil, 0, err
		}

		returns = append(returns, returnReq)
	}

	return returns, total, nil
}

func (r *returnRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.ReturnStatus, refundAmount float64) error {
	query := `
        UPDATE returns
        SET status = $1, refund_amount = $2, updated_at = NOW()
        WHERE id = $3
    `

	_, err := r.db.Exec(ctx, query, status, refundAmount, id)
	return err
}

func (r *returnRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]models.Return, error) {
	query := `
        SELECT id, order_id, user_id, reason, status, refund_amount, created_at, updated_at
        FROM returns
        WHERE order_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var returns []models.Return
	for rows.Next() {
		var returnReq models.Return
		err := rows.Scan(
			&returnReq.ID,
			&returnReq.OrderID,
			&returnReq.UserID,
			&returnReq.Reason,
			&returnReq.Status,
			&returnReq.RefundAmount,
			&returnReq.CreatedAt,
			&returnReq.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		returns = append(returns, returnReq)
	}

	return returns, nil
}

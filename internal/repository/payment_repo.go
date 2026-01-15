package repository

import (
	"context"

	"ecommerce-backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *models.Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Payment, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Payment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.PaymentStatus, transactionID string) error
	UpdateStatusWithRefund(ctx context.Context, id uuid.UUID, status models.PaymentStatus, refundAmount float64) error
}

type paymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	query := `
        INSERT INTO payments (id, order_id, amount, status, payment_method, transaction_id, payment_details)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING created_at, updated_at
    `

	return r.db.QueryRow(ctx, query,
		payment.ID,
		payment.OrderID,
		payment.Amount,
		payment.Status,
		payment.PaymentMethod,
		payment.TransactionID,
		payment.PaymentDetails,
	).Scan(&payment.CreatedAt, &payment.UpdatedAt)
}

func (r *paymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	query := `
        SELECT id, order_id, amount, status, payment_method, transaction_id, 
               payment_details, created_at, updated_at
        FROM payments
        WHERE id = $1
    `

	var payment models.Payment
	err := r.db.QueryRow(ctx, query, id).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Amount,
		&payment.Status,
		&payment.PaymentMethod,
		&payment.TransactionID,
		&payment.PaymentDetails,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *paymentRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Payment, error) {
	query := `
        SELECT id, order_id, amount, status, payment_method, transaction_id, 
               payment_details, created_at, updated_at
        FROM payments
        WHERE order_id = $1
        ORDER BY created_at DESC
        LIMIT 1
    `

	var payment models.Payment
	err := r.db.QueryRow(ctx, query, orderID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Amount,
		&payment.Status,
		&payment.PaymentMethod,
		&payment.TransactionID,
		&payment.PaymentDetails,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.PaymentStatus, transactionID string) error {
	query := `
        UPDATE payments
        SET status = $1, transaction_id = $2, updated_at = NOW()
        WHERE id = $3
    `

	_, err := r.db.Exec(ctx, query, status, transactionID, id)
	return err
}

func (r *paymentRepository) UpdateStatusWithRefund(ctx context.Context, id uuid.UUID, status models.PaymentStatus, refundAmount float64) error {
	query := `
        UPDATE payments
        SET status = $1, updated_at = NOW()
        WHERE id = $2
    `

	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

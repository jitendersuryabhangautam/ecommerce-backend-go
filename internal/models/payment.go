package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentPending    PaymentStatus = "pending"
	PaymentProcessing PaymentStatus = "processing"
	PaymentCompleted  PaymentStatus = "completed"
	PaymentFailed     PaymentStatus = "failed"
	PaymentRefunded   PaymentStatus = "refunded"
)

type Payment struct {
	ID             uuid.UUID              `json:"id"`
	OrderID        uuid.UUID              `json:"order_id"`
	Amount         float64                `json:"amount"`
	Status         PaymentStatus          `json:"status"`
	PaymentMethod  string                 `json:"payment_method"`
	TransactionID  string                 `json:"transaction_id"`
	PaymentDetails map[string]interface{} `json:"payment_details"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

type CreatePaymentRequest struct {
	OrderID       uuid.UUID `json:"order_id" validate:"required"`
	PaymentMethod string    `json:"payment_method" validate:"required"`
}

type VerifyPaymentRequest struct {
	PaymentID     uuid.UUID `json:"payment_id" validate:"required"`
	TransactionID string    `json:"transaction_id" validate:"required"`
}

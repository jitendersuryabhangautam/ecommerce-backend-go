package models

import (
	"time"

	"github.com/google/uuid"
)

type ReturnStatus string

const (
	ReturnRequested ReturnStatus = "requested"
	ReturnApproved  ReturnStatus = "approved"
	ReturnRejected  ReturnStatus = "rejected"
	ReturnCompleted ReturnStatus = "completed"
)

type Return struct {
	ID           uuid.UUID    `json:"id"`
	OrderID      uuid.UUID    `json:"order_id"`
	UserID       uuid.UUID    `json:"user_id"`
	Reason       string       `json:"reason"`
	Status       ReturnStatus `json:"status"`
	RefundAmount float64      `json:"refund_amount"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type CreateReturnRequest struct {
	OrderID uuid.UUID `json:"order_id" validate:"required"`
	Reason  string    `json:"reason" validate:"required"`
}

type ProcessReturnRequest struct {
	Status       ReturnStatus `json:"status" validate:"required"`
	RefundAmount float64      `json:"refund_amount"`
}

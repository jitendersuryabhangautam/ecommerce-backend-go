package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderPending    OrderStatus = "pending"
	OrderProcessing OrderStatus = "processing"
	OrderShipped    OrderStatus = "shipped"
	OrderDelivered  OrderStatus = "delivered"
	OrderCompleted  OrderStatus = "completed"
	OrderCancelled  OrderStatus = "cancelled"
	OrderRefunded   OrderStatus = "refunded"
	OrderReturnRequested OrderStatus = "return_requested"
)

type Order struct {
	ID              uuid.UUID   `json:"id"`
	UserID          uuid.UUID   `json:"user_id"`
	OrderNumber     string      `json:"order_number"`
	TotalAmount     float64     `json:"total_amount"`
	Status          OrderStatus `json:"status"`
	PaymentMethod   string      `json:"payment_method"`
	ShippingAddress Address     `json:"shipping_address"`
	BillingAddress  Address     `json:"billing_address"`
	Items           []OrderItem `json:"items"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID          uuid.UUID `json:"id"`
	OrderID     uuid.UUID `json:"order_id"`
	ProductID   uuid.UUID `json:"product_id"`
	Product     Product   `json:"product"`
	Quantity    int       `json:"quantity"`
	PriceAtTime float64   `json:"price_at_time"`
	CreatedAt   time.Time `json:"created_at"`
}

type Address struct {
	FullName   string `json:"full_name"`
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
	Phone      string `json:"phone"`
}

type CreateOrderRequest struct {
	ShippingAddress Address `json:"shipping_address" validate:"required"`
	BillingAddress  Address `json:"billing_address" validate:"required"`
	PaymentMethod   string  `json:"payment_method" validate:"required,oneof=cc dc cod"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" validate:"required"`
}

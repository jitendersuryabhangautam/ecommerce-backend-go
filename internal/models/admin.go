package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminUserSummary struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type AdminOrderItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

type AdminOrder struct {
	ID              uuid.UUID        `json:"id"`
	UserID          uuid.UUID        `json:"user_id"`
	User            AdminUserSummary `json:"user"`
	OrderNumber     string           `json:"order_number"`
	TotalAmount     float64          `json:"total_amount"`
	Status          OrderStatus      `json:"status"`
	PaymentMethod   string           `json:"payment_method"`
	ShippingAddress Address          `json:"shipping_address"`
	BillingAddress  Address          `json:"billing_address"`
	Items           []AdminOrderItem `json:"items"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

type AdminReturnOrderSummary struct {
	OrderNumber string `json:"order_number"`
}

type AdminReturn struct {
	ID           uuid.UUID              `json:"id"`
	OrderID      uuid.UUID              `json:"order_id"`
	Order        AdminReturnOrderSummary `json:"order"`
	User         AdminUserSummary       `json:"user"`
	Reason       string                 `json:"reason"`
	Status       ReturnStatus           `json:"status"`
	RefundAmount float64                `json:"refund_amount"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type AdminTotals struct {
	TotalRevenue   float64 `json:"total_revenue"`
	TotalOrders    int     `json:"total_orders"`
	TotalProducts  int     `json:"total_products"`
	TotalCustomers int     `json:"total_customers"`
	AvgOrderValue  float64 `json:"avg_order_value"`
}

type AdminStatusCount struct {
	Status OrderStatus `json:"status"`
	Count  int         `json:"count"`
}

type AdminAnalytics struct {
	RangeDays      int               `json:"range_days"`
	Totals         AdminTotals       `json:"totals"`
	OrdersByStatus []AdminStatusCount `json:"orders_by_status"`
}

type TopProductItem struct {
	Product       Product `json:"product"`
	TotalQuantity int     `json:"total_quantity"`
	TotalRevenue  float64 `json:"total_revenue"`
}

type TopProductsResponse struct {
	RangeDays int              `json:"range_days"`
	Items     []TopProductItem `json:"items"`
}

type RecentOrdersResponse struct {
	Limit  int          `json:"limit"`
	Orders []AdminOrder `json:"orders"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin customer"`
}

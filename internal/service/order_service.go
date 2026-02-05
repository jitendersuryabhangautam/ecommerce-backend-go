package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repository"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, req models.CreateOrderRequest) (*models.Order, error)
	GetOrder(ctx context.Context, orderID, userID uuid.UUID) (*models.Order, error)
	GetUserOrders(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Order, int, error)
	GetAllOrders(ctx context.Context, page, limit int, status string, rangeDays int) ([]models.AdminOrder, int, error)
	GetOrderAdmin(ctx context.Context, orderID uuid.UUID) (*models.AdminOrder, error)
	GetRecentOrders(ctx context.Context, limit, rangeDays int) ([]models.AdminOrder, error)
	GetAnalytics(ctx context.Context, rangeDays int) (*models.AdminAnalytics, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error
	CancelOrder(ctx context.Context, orderID, userID uuid.UUID) error
	ProcessOrderReturn(ctx context.Context, orderID uuid.UUID, returnID uuid.UUID) error
}

type orderService struct {
	orderRepo   repository.OrderRepository
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	cartSvc     CartService
	paymentSvc  PaymentService
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	cartRepo repository.CartRepository,
	productRepo repository.ProductRepository,
	cartSvc CartService,
	paymentSvc PaymentService,
) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		cartSvc:     cartSvc,
		paymentSvc:  paymentSvc,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, userID uuid.UUID, req models.CreateOrderRequest) (*models.Order, error) {
	// Get user's cart
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Validate cart
	valid, validationErrors, err := s.cartSvc.ValidateCart(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, fmt.Errorf("cart validation failed: %v", validationErrors)
	}

	// Start transaction
	tx, err := s.orderRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Calculate total and prepare order items
	var totalAmount float64
	var orderItems []models.OrderItem

	for _, cartItem := range cart.Items {
		// Calculate item total
		itemTotal := cartItem.Product.Price * float64(cartItem.Quantity)
		totalAmount += itemTotal

		// Prepare order item
		orderItem := models.OrderItem{
			ID:          uuid.New(),
			ProductID:   cartItem.ProductID,
			Product:     cartItem.Product,
			Quantity:    cartItem.Quantity,
			PriceAtTime: cartItem.Product.Price,
			CreatedAt:   time.Now(),
		}
		orderItems = append(orderItems, orderItem)

		// Deduct stock from inventory (within transaction)
		err = s.productRepo.UpdateStockWithTx(ctx, tx, cartItem.ProductID, -cartItem.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to update stock for product %s: %w",
				cartItem.ProductID, err)
		}
	}

	// Create order
	order := &models.Order{
		ID:              uuid.New(),
		UserID:          userID,
		OrderNumber:     generateOrderNumber(),
		TotalAmount:     totalAmount,
		Status:          models.OrderPending,
		PaymentMethod:   req.PaymentMethod,
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
		Items:           orderItems,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Create order in transaction
	err = s.orderRepo.CreateWithTx(ctx, tx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Clear cart
	err = s.cartSvc.ClearCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to clear cart: %w", err)
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Create payment immediately for card payments
	if req.PaymentMethod == "cc" || req.PaymentMethod == "dc" {
		_, err = s.paymentSvc.CreatePaymentForOrder(ctx, order.ID, req.PaymentMethod, models.PaymentCompleted)
		if err != nil {
			return nil, err
		}
	}

	return order, nil
}

func generateOrderNumber() string {
	timestamp := time.Now().Unix()
	random := uuid.New().String()[:8]
	return fmt.Sprintf("ORD-%d-%s", timestamp, random)
}

func (s *orderService) GetOrder(ctx context.Context, orderID, userID uuid.UUID) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.New("order not found")
	}

	// Check if user is authorized to view this order
	if order.UserID != userID {
		return nil, errors.New("unauthorized to view this order")
	}

	return order, nil
}

func (s *orderService) GetUserOrders(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Order, int, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 50 {
		limit = 10
	}

	return s.orderRepo.GetByUserID(ctx, userID, page, limit)
}

func (s *orderService) GetAllOrders(ctx context.Context, page, limit int, status string, rangeDays int) ([]models.AdminOrder, int, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 50 {
		limit = 10
	}

	return s.orderRepo.GetAll(ctx, page, limit, status, rangeDays)
}

func (s *orderService) GetOrderAdmin(ctx context.Context, orderID uuid.UUID) (*models.AdminOrder, error) {
	fmt.Printf("[ORDER SERVICE] GetOrderAdmin called for orderID: %s\n", orderID.String())
	order, err := s.orderRepo.GetAdminByID(ctx, orderID)
	if err != nil {
		fmt.Printf("[ORDER SERVICE ERROR] Repository error: %v\n", err)
		return nil, err
	}
	if order == nil {
		fmt.Printf("[ORDER SERVICE] Order not found in repository\n")
		return nil, errors.New("order not found")
	}
	fmt.Printf("[ORDER SERVICE SUCCESS] Order found: %s\n", order.OrderNumber)
	return order, nil
}

func (s *orderService) GetRecentOrders(ctx context.Context, limit, rangeDays int) ([]models.AdminOrder, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.orderRepo.GetRecent(ctx, limit, rangeDays)
}

func (s *orderService) GetAnalytics(ctx context.Context, rangeDays int) (*models.AdminAnalytics, error) {
	return s.orderRepo.GetAnalytics(ctx, rangeDays)
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error {
	// Check if order exists
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("order not found")
	}

	// If status is already the same, no update needed
	if order.Status == status {
		return nil
	}

	// Validate status transition
	if !isValidStatusTransition(order.Status, status) {
		return fmt.Errorf("invalid status transition from %s to %s", order.Status, status)
	}

	if err := s.orderRepo.UpdateStatus(ctx, orderID, status); err != nil {
		return err
	}

	// For COD, create payment when delivered
	if status == models.OrderDelivered && order.PaymentMethod == "cod" {
		existing, err := s.paymentSvc.GetPaymentByOrderID(ctx, orderID)
		if err != nil {
			return err
		}
		if existing == nil {
			_, err := s.paymentSvc.CreatePaymentForOrder(ctx, orderID, "cod", models.PaymentCompleted)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isValidStatusTransition(from, to models.OrderStatus) bool {
	transitions := map[models.OrderStatus][]models.OrderStatus{
		models.OrderPending:    {models.OrderProcessing, models.OrderCancelled},
		models.OrderProcessing: {models.OrderShipped, models.OrderCancelled},
		models.OrderShipped:    {models.OrderDelivered},
		models.OrderDelivered:  {models.OrderCompleted},
		models.OrderCompleted:  {},
		models.OrderCancelled:  {},
		models.OrderRefunded:   {},
	}

	allowed, ok := transitions[from]
	if !ok {
		return false
	}

	for _, s := range allowed {
		if s == to {
			return true
		}
	}

	return false
}

func (s *orderService) CancelOrder(ctx context.Context, orderID, userID uuid.UUID) error {
	// Get order and verify ownership
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("order not found")
	}

	if order.UserID != userID {
		return errors.New("unauthorized to cancel this order")
	}

	// Check if order can be cancelled
	if order.Status != models.OrderPending && order.Status != models.OrderProcessing {
		return errors.New("order cannot be cancelled at this stage")
	}

	// Start transaction to cancel order and restore stock
	tx, err := s.orderRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Restore stock for each item
	for _, item := range order.Items {
		err = s.productRepo.UpdateStock(ctx, item.ProductID, item.Quantity)
		if err != nil {
			return fmt.Errorf("failed to restore stock for product %s: %w",
				item.ProductID, err)
		}
	}

	// Update order status
	err = s.orderRepo.UpdateStatus(ctx, orderID, models.OrderCancelled)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *orderService) ProcessOrderReturn(ctx context.Context, orderID uuid.UUID, returnID uuid.UUID) error {
	// This would integrate with the return service
	// For now, just update order status to refunded
	return s.orderRepo.UpdateStatus(ctx, orderID, models.OrderRefunded)
}

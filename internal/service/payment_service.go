package service

import (
	"context"
	"errors"
	"time"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repository"

	"github.com/google/uuid"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, req models.CreatePaymentRequest, userID uuid.UUID) (*models.Payment, error)
	VerifyPayment(ctx context.Context, req models.VerifyPaymentRequest) (*models.Payment, error)
	ProcessRefund(ctx context.Context, paymentID uuid.UUID, amount float64) error
	GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Payment, error)
	CreatePaymentForOrder(ctx context.Context, orderID uuid.UUID, method string, status models.PaymentStatus) (*models.Payment, error)
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
	orderRepo   repository.OrderRepository
}

func NewPaymentService(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
	}
}

func (s *paymentService) CreatePayment(ctx context.Context, req models.CreatePaymentRequest, userID uuid.UUID) (*models.Payment, error) {
	// Get order
	order, err := s.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.New("order not found")
	}

	// Verify order belongs to user
	if order.UserID != userID {
		return nil, errors.New("unauthorized to create payment for this order")
	}

	// Check if payment already exists
	existingPayment, err := s.paymentRepo.GetByOrderID(ctx, req.OrderID)
	if err == nil && existingPayment != nil {
		return nil, errors.New("payment already exists for this order")
	}

	// Create payment
	transactionID := "TXN-" + uuid.New().String()[:8]
	payment := &models.Payment{
		ID:             uuid.New(),
		OrderID:        req.OrderID,
		Amount:         order.TotalAmount,
		Status:         models.PaymentPending,
		PaymentMethod:  req.PaymentMethod,
		TransactionID:  transactionID,
		PaymentDetails: make(map[string]interface{}),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = s.paymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, err
	}

	// For demo purposes, simulate payment processing
	go s.simulatePaymentProcessing(ctx, payment.ID)

	return payment, nil
}

func (s *paymentService) simulatePaymentProcessing(ctx context.Context, paymentID uuid.UUID) {
	// Simulate payment processing delay
	time.Sleep(3 * time.Second)

	// Simulate successful payment
	transactionID := "TXN-" + uuid.New().String()[:8]
	s.paymentRepo.UpdateStatus(ctx, paymentID, models.PaymentCompleted, transactionID)

	// Update order status
	payment, _ := s.paymentRepo.GetByID(ctx, paymentID)
	if payment != nil {
		s.orderRepo.UpdateStatus(ctx, payment.OrderID, models.OrderProcessing)
	}
}

func (s *paymentService) VerifyPayment(ctx context.Context, req models.VerifyPaymentRequest) (*models.Payment, error) {
	payment, err := s.paymentRepo.GetByID(ctx, req.PaymentID)
	if err != nil {
		return nil, err
	}

	if payment == nil {
		return nil, errors.New("payment not found")
	}

	// Verify transaction
	if payment.TransactionID != req.TransactionID {
		return nil, errors.New("invalid transaction ID")
	}

	return payment, nil
}

func (s *paymentService) ProcessRefund(ctx context.Context, paymentID uuid.UUID, amount float64) error {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return err
	}

	if payment == nil {
		return errors.New("payment not found")
	}

	if payment.Status != models.PaymentCompleted {
		return errors.New("can only refund completed payments")
	}

	if amount > payment.Amount {
		return errors.New("refund amount cannot exceed payment amount")
	}

	// Update payment status
	err = s.paymentRepo.UpdateStatusWithRefund(ctx, paymentID, models.PaymentRefunded, amount)
	if err != nil {
		return err
	}

	// Update order status
	return s.orderRepo.UpdateStatus(ctx, payment.OrderID, models.OrderRefunded)
}

func (s *paymentService) GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Payment, error) {
	return s.paymentRepo.GetByOrderID(ctx, orderID)
}

func (s *paymentService) CreatePaymentForOrder(ctx context.Context, orderID uuid.UUID, method string, status models.PaymentStatus) (*models.Payment, error) {
	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	// Check if payment already exists
	existingPayment, err := s.paymentRepo.GetByOrderID(ctx, orderID)
	if err == nil && existingPayment != nil {
		return nil, errors.New("payment already exists for this order")
	}

	transactionID := "TXN-" + uuid.New().String()[:8]
	payment := &models.Payment{
		ID:             uuid.New(),
		OrderID:        orderID,
		Amount:         order.TotalAmount,
		Status:         status,
		PaymentMethod:  method,
		TransactionID:  transactionID,
		PaymentDetails: make(map[string]interface{}),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	if status == models.PaymentCompleted {
		_ = s.orderRepo.UpdateStatus(ctx, orderID, models.OrderProcessing)
	}

	return payment, nil
}

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

type ReturnService interface {
	CreateReturn(ctx context.Context, req models.CreateReturnRequest, userID uuid.UUID) (*models.Return, error)
	GetReturn(ctx context.Context, returnID uuid.UUID, userID uuid.UUID) (*models.Return, error)
	GetUserReturns(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Return, int, error)
	GetAllReturns(ctx context.Context, page, limit int, status string, rangeDays int) ([]models.AdminReturn, int, error)
	ProcessReturn(ctx context.Context, returnID uuid.UUID, req models.ProcessReturnRequest) (*models.Return, error)
}

type returnService struct {
	returnRepo  repository.ReturnRepository
	orderRepo   repository.OrderRepository
	paymentSvc  PaymentService
	productRepo repository.ProductRepository
}

func NewReturnService(
	returnRepo repository.ReturnRepository,
	orderRepo repository.OrderRepository,
	paymentSvc PaymentService,
	productRepo repository.ProductRepository,
) ReturnService {
	return &returnService{
		returnRepo:  returnRepo,
		orderRepo:   orderRepo,
		paymentSvc:  paymentSvc,
		productRepo: productRepo,
	}
}

func (s *returnService) CreateReturn(ctx context.Context, req models.CreateReturnRequest, userID uuid.UUID) (*models.Return, error) {
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
		return nil, errors.New("unauthorized to create return for this order")
	}

	// Check if order can be returned
	if order.Status != models.OrderDelivered && order.Status != models.OrderCompleted {
		return nil, errors.New("order cannot be returned at this stage")
	}

	// Check if return period has expired (14 days from delivery)
	deliveryTime := order.CreatedAt.Add(7 * 24 * time.Hour) // Assuming 7 days for delivery
	if time.Since(deliveryTime) > 14*24*time.Hour {
		return nil, errors.New("return period has expired")
	}

	// Check if return already exists for this order
	existingReturns, err := s.returnRepo.GetByOrderID(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	for _, r := range existingReturns {
		if r.Status == models.ReturnRequested || r.Status == models.ReturnApproved {
			return nil, errors.New("return already requested for this order")
		}
	}

	// Create return request
	returnReq := &models.Return{
		ID:           uuid.New(),
		OrderID:      req.OrderID,
		UserID:       userID,
		Reason:       req.Reason,
		Status:       models.ReturnRequested,
		RefundAmount: order.TotalAmount,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.returnRepo.Create(ctx, returnReq)
	if err != nil {
		return nil, err
	}

	// Mark order as return requested
	if err := s.orderRepo.UpdateStatus(ctx, order.ID, models.OrderReturnRequested); err != nil {
		return nil, err
	}

	return returnReq, nil
}

func (s *returnService) GetReturn(ctx context.Context, returnID uuid.UUID, userID uuid.UUID) (*models.Return, error) {
	returnReq, err := s.returnRepo.GetByID(ctx, returnID)
	if err != nil {
		return nil, err
	}

	if returnReq == nil {
		return nil, errors.New("return not found")
	}

	// Check if user is authorized to view this return
	if returnReq.UserID != userID {
		return nil, errors.New("unauthorized to view this return")
	}

	return returnReq, nil
}

func (s *returnService) GetUserReturns(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Return, int, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 50 {
		limit = 10
	}

	return s.returnRepo.GetByUserID(ctx, userID, page, limit)
}

func (s *returnService) GetAllReturns(ctx context.Context, page, limit int, status string, rangeDays int) ([]models.AdminReturn, int, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 50 {
		limit = 10
	}

	return s.returnRepo.GetAll(ctx, page, limit, status, rangeDays)
}

func (s *returnService) ProcessReturn(ctx context.Context, returnID uuid.UUID, req models.ProcessReturnRequest) (*models.Return, error) {
	returnReq, err := s.returnRepo.GetByID(ctx, returnID)
	if err != nil {
		return nil, err
	}

	if returnReq == nil {
		return nil, errors.New("return not found")
	}

	// Get order
	order, err := s.orderRepo.GetByID(ctx, returnReq.OrderID)
	if err != nil {
		return nil, err
	}

	// Process based on status
	if req.Status == models.ReturnApproved {
		// Calculate refund amount (full refund for demo)
		refundAmount := req.RefundAmount
		if refundAmount == 0 {
			refundAmount = order.TotalAmount
		}

		// Process refund through payment service
		payment, err := s.paymentSvc.GetPaymentByOrderID(ctx, returnReq.OrderID)
		if err != nil {
			return nil, err
		}

		if payment != nil {
			err = s.paymentSvc.ProcessRefund(ctx, payment.ID, refundAmount)
			if err != nil {
				return nil, err
			}
		}

		// Restore stock for order items
		for _, item := range order.Items {
			err = s.productRepo.UpdateStock(ctx, item.ProductID, item.Quantity)
			if err != nil {
				return nil, fmt.Errorf("failed to restore stock for product %s: %w",
					item.ProductID, err)
			}
		}

		// Update return with refund amount
		returnReq.RefundAmount = refundAmount
	}

	// Update return status
	err = s.returnRepo.UpdateStatus(ctx, returnID, req.Status, returnReq.RefundAmount)
	if err != nil {
		return nil, err
	}

	// Get updated return
	return s.returnRepo.GetByID(ctx, returnID)
}

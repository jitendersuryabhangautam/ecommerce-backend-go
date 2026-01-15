package service

import (
	"context"
	"errors"
	"fmt"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repository"

	"github.com/google/uuid"
)

type CartService interface {
	GetCart(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
	AddToCart(ctx context.Context, userID uuid.UUID, req models.AddToCartRequest) (*models.Cart, error)
	UpdateCartItem(ctx context.Context, userID, itemID uuid.UUID, req models.UpdateCartItemRequest) (*models.Cart, error)
	RemoveFromCart(ctx context.Context, userID, itemID uuid.UUID) (*models.Cart, error)
	ClearCart(ctx context.Context, userID uuid.UUID) error
	ValidateCart(ctx context.Context, cartID uuid.UUID) (bool, []string, error)
}

type cartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	productSvc  ProductService
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository, productSvc ProductService) CartService {
	return &cartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
		productSvc:  productSvc,
	}
}

func (s *cartService) GetCart(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *cartService) AddToCart(ctx context.Context, userID uuid.UUID, req models.AddToCartRequest) (*models.Cart, error) {
	// Get or create cart
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check product exists and has enough stock
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New("product not found")
	}

	// Check available stock
	available, err := s.productSvc.CheckStock(ctx, req.ProductID, req.Quantity)
	if err != nil {
		return nil, err
	}

	if !available {
		return nil, errors.New("insufficient stock")
	}

	// Reserve stock
	err = s.productSvc.ReserveStock(ctx, req.ProductID, cart.ID, req.Quantity)
	if err != nil {
		return nil, err
	}

	// Add to cart
	err = s.cartRepo.AddItem(ctx, cart.ID, req.ProductID, req.Quantity)
	if err != nil {
		// Release reservation if adding to cart fails
		s.productSvc.ReleaseStockReservation(ctx, req.ProductID, cart.ID)
		return nil, err
	}

	// Get updated cart
	return s.cartRepo.GetCartWithItems(ctx, cart.ID)
}

func (s *cartService) UpdateCartItem(ctx context.Context, userID, itemID uuid.UUID, req models.UpdateCartItemRequest) (*models.Cart, error) {
	// Get cart
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Find the item to update
	var itemToUpdate *models.CartItem
	for _, item := range cart.Items {
		if item.ID == itemID {
			itemToUpdate = &item
			break
		}
	}

	if itemToUpdate == nil {
		return nil, errors.New("cart item not found")
	}

	// Calculate quantity difference
	quantityDiff := req.Quantity - itemToUpdate.Quantity

	if quantityDiff > 0 {
		// Need more stock - check availability
		available, err := s.productSvc.CheckStock(ctx, itemToUpdate.ProductID, quantityDiff)
		if err != nil {
			return nil, err
		}

		if !available {
			return nil, errors.New("insufficient stock for additional quantity")
		}

		// Reserve additional stock
		err = s.productSvc.ReserveStock(ctx, itemToUpdate.ProductID, cart.ID, quantityDiff)
		if err != nil {
			return nil, err
		}
	} else if quantityDiff < 0 {
		// Releasing stock
		err = s.productSvc.ReleaseStockReservation(ctx, itemToUpdate.ProductID, cart.ID)
		if err != nil {
			return nil, err
		}
	}

	// Update cart item
	err = s.cartRepo.UpdateItem(ctx, cart.ID, itemID, req.Quantity)
	if err != nil {
		return nil, err
	}

	// Get updated cart
	return s.cartRepo.GetCartWithItems(ctx, cart.ID)
}

func (s *cartService) RemoveFromCart(ctx context.Context, userID, itemID uuid.UUID) (*models.Cart, error) {
	// Get cart
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Find the item to remove
	var itemToRemove *models.CartItem
	for _, item := range cart.Items {
		if item.ID == itemID {
			itemToRemove = &item
			break
		}
	}

	if itemToRemove == nil {
		return nil, errors.New("cart item not found")
	}

	// Release stock reservation
	err = s.productSvc.ReleaseStockReservation(ctx, itemToRemove.ProductID, cart.ID)
	if err != nil {
		return nil, err
	}

	// Remove from cart
	err = s.cartRepo.RemoveItem(ctx, cart.ID, itemID)
	if err != nil {
		return nil, err
	}

	// Get updated cart
	return s.cartRepo.GetCartWithItems(ctx, cart.ID)
}

func (s *cartService) ClearCart(ctx context.Context, userID uuid.UUID) error {
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Release all stock reservations
	for _, item := range cart.Items {
		s.productSvc.ReleaseStockReservation(ctx, item.ProductID, cart.ID)
	}

	// Clear cart
	return s.cartRepo.ClearCart(ctx, cart.ID)
}

func (s *cartService) ValidateCart(ctx context.Context, cartID uuid.UUID) (bool, []string, error) {
	cart, err := s.cartRepo.GetCartWithItems(ctx, cartID)
	if err != nil {
		return false, nil, err
	}

	var errors []string
	valid := true

	for _, item := range cart.Items {
		available, err := s.productSvc.CheckStock(ctx, item.ProductID, item.Quantity)
		if err != nil {
			return false, nil, err
		}

		if !available {
			valid = false
			errors = append(errors,
				fmt.Sprintf("Insufficient stock for %s. Available: %d",
					item.Product.Name, item.Product.Stock))
		}
	}

	return valid, errors, nil
}

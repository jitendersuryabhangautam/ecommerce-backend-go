package service

import (
	"context"
	"errors"
	"time"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repository"

	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req models.ProductRequest) (*models.Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetProducts(ctx context.Context, page, limit int, category, search string) ([]models.Product, int, error)
	GetAdminProducts(ctx context.Context, page, limit, rangeDays int) ([]models.Product, int, error)
	GetTopProducts(ctx context.Context, limit, rangeDays int) ([]models.TopProductItem, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req models.ProductUpdateRequest) (*models.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	CheckStock(ctx context.Context, productID uuid.UUID, quantity int) (bool, error)
	ReserveStock(ctx context.Context, productID, cartID uuid.UUID, quantity int) error
	ReleaseStockReservation(ctx context.Context, productID, cartID uuid.UUID) error
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) CreateProduct(ctx context.Context, req models.ProductRequest) (*models.Product, error) {
	// Check if SKU already exists
	existingProduct, err := s.productRepo.GetBySKU(ctx, req.SKU)
	if err != nil {
		return nil, err
	}

	if existingProduct != nil {
		return nil, errors.New("product with this SKU already exists")
	}

	product := &models.Product{
		ID:          uuid.New(),
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
	}

	err = s.productRepo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New("product not found")
	}

	return product, nil
}

func (s *productService) GetProducts(ctx context.Context, page, limit int, category, search string) ([]models.Product, int, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.productRepo.GetAll(ctx, page, limit, category, search)
}

func (s *productService) GetAdminProducts(ctx context.Context, page, limit, rangeDays int) ([]models.Product, int, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.productRepo.GetAllAdmin(ctx, page, limit, rangeDays)
}

func (s *productService) GetTopProducts(ctx context.Context, limit, rangeDays int) ([]models.TopProductItem, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.productRepo.GetTopProducts(ctx, limit, rangeDays)
}

func (s *productService) UpdateProduct(ctx context.Context, id uuid.UUID, req models.ProductUpdateRequest) (*models.Product, error) {
	// Check if product exists
	existingProduct, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingProduct == nil {
		return nil, errors.New("product not found")
	}

	err = s.productRepo.Update(ctx, id, &req)
	if err != nil {
		return nil, err
	}

	// Get updated product
	return s.productRepo.GetByID(ctx, id)
}

func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	// Check if product exists
	existingProduct, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingProduct == nil {
		return errors.New("product not found")
	}

	return s.productRepo.Delete(ctx, id)
}

func (s *productService) CheckStock(ctx context.Context, productID uuid.UUID, quantity int) (bool, error) {
	available, err := s.productRepo.GetAvailableStock(ctx, productID)
	if err != nil {
		return false, err
	}

	return available >= quantity, nil
}

func (s *productService) ReserveStock(ctx context.Context, productID, cartID uuid.UUID, quantity int) error {
	expiresAt := time.Now().Add(10 * time.Minute).Unix()
	return s.productRepo.ReserveStock(ctx, productID, cartID, quantity, expiresAt)
}

func (s *productService) ReleaseStockReservation(ctx context.Context, productID, cartID uuid.UUID) error {
	return s.productRepo.ReleaseStockReservation(ctx, productID, cartID)
}

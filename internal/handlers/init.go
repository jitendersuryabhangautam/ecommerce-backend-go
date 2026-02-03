package handlers

import (
	"ecommerce-backend/internal/config"
	"ecommerce-backend/internal/repository"
	"ecommerce-backend/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	AuthHandler    *AuthHandler
	ProductHandler *ProductHandler
	CartHandler    *CartHandler
	OrderHandler   *OrderHandler
	PaymentHandler *PaymentHandler
	ReturnHandler  *ReturnHandler
	HealthHandler  *HealthHandler
}

func InitRepositories(db *pgxpool.Pool, cfg *config.Config) *Repositories {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	returnRepo := repository.NewReturnRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiry)
	productService := service.NewProductService(productRepo)
	cartService := service.NewCartService(cartRepo, productRepo, productService)
	paymentService := service.NewPaymentService(paymentRepo, orderRepo)
	orderService := service.NewOrderService(orderRepo, cartRepo, productRepo, cartService, paymentService)
	returnService := service.NewReturnService(returnRepo, orderRepo, paymentService, productRepo)

	// Initialize handlers
	authHandler := NewAuthHandler(authService)
	productHandler := NewProductHandler(productService)
	cartHandler := NewCartHandler(cartService)
	orderHandler := NewOrderHandler(orderService)
	paymentHandler := NewPaymentHandler(paymentService)
	returnHandler := NewReturnHandler(returnService)
	healthHandler := NewHealthHandler(db)

	return &Repositories{
		AuthHandler:    authHandler,
		ProductHandler: productHandler,
		CartHandler:    cartHandler,
		OrderHandler:   orderHandler,
		PaymentHandler: paymentHandler,
		ReturnHandler:  returnHandler,
		HealthHandler:  healthHandler,
	}
}

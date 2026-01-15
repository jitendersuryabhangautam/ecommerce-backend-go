package main

import (
	"log"

	"ecommerce-backend/internal/config"
	"ecommerce-backend/internal/handlers"
	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer db.Close()

	log.Println("✅ Database connection established")

	// Set Gin mode
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.GinCORSMiddleware())
	router.Use(middleware.GinRecovery())
	router.Use(middleware.GinLogging())
	router.Use(middleware.GinRequestID())

	// Initialize repositories, services, and handlers
	repos := handlers.InitRepositories(db, cfg)

	// Health check endpoints (public)
	router.GET("/health", repos.HealthHandler.HealthCheck)
	router.GET("/ready", repos.HealthHandler.ReadinessCheck)
	router.GET("/metrics", repos.HealthHandler.Metrics)

	// API version prefix
	api := router.Group("/api/v1")

	// Public routes
	{
		// Auth routes
		api.POST("/auth/register", repos.AuthHandler.Register)
		api.POST("/auth/login", repos.AuthHandler.Login)

		// Product routes (public read access)
		api.GET("/products", repos.ProductHandler.GetProducts)
		api.GET("/products/:id", repos.ProductHandler.GetProduct)
	}

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(middleware.GinAuthMiddleware(repos.AuthHandler.AuthService))
	{
		// User routes
		protected.GET("/users/profile", repos.AuthHandler.GetProfile)
		protected.PUT("/users/profile", repos.AuthHandler.UpdateProfile)
		protected.PUT("/users/change-password", repos.AuthHandler.ChangePassword)

		// Cart routes
		protected.GET("/cart", repos.CartHandler.GetCart)
		protected.GET("/cart/validate", repos.CartHandler.ValidateCart)
		protected.POST("/cart/items", repos.CartHandler.AddToCart)
		protected.PUT("/cart/items/:itemId", repos.CartHandler.UpdateCartItem)
		protected.DELETE("/cart/items/:itemId", repos.CartHandler.RemoveFromCart)
		protected.DELETE("/cart", repos.CartHandler.ClearCart)

		// Order routes
		protected.POST("/orders", repos.OrderHandler.CreateOrder)
		protected.GET("/orders", repos.OrderHandler.GetUserOrders)
		protected.GET("/orders/:id/payment", repos.PaymentHandler.GetPaymentByOrder)
		protected.GET("/orders/:id", repos.OrderHandler.GetOrder)
		protected.PUT("/orders/:id/cancel", repos.OrderHandler.CancelOrder)

		// Payment routes
		protected.POST("/payments", repos.PaymentHandler.CreatePayment)
		protected.POST("/payments/:id/verify", repos.PaymentHandler.VerifyPayment)

		// Return routes
		protected.POST("/returns", repos.ReturnHandler.CreateReturn)
		protected.GET("/returns", repos.ReturnHandler.GetUserReturns)
		protected.GET("/returns/:id", repos.ReturnHandler.GetReturn)
	}

	// Admin routes (require admin role)
	admin := api.Group("/admin")
	admin.Use(middleware.GinAuthMiddleware(repos.AuthHandler.AuthService))
	admin.Use(middleware.GinAdminMiddleware())
	{
		// Product management
		admin.POST("/products", repos.ProductHandler.CreateProduct)
		admin.PUT("/products/:id", repos.ProductHandler.UpdateProduct)
		admin.DELETE("/products/:id", repos.ProductHandler.DeleteProduct)

		// Order management
		admin.GET("/orders", repos.OrderHandler.GetAllOrders)
		admin.PUT("/orders/:id/status", repos.OrderHandler.UpdateOrderStatus)

		// User management
		admin.GET("/users", repos.AuthHandler.GetAllUsers)
		admin.PUT("/users/:id/role", repos.AuthHandler.UpdateUserRole)
	}

	// Print API documentation
	log.Println("📚 API Documentation available at http://localhost:" + cfg.Port)
	log.Println("🚀 Environment: " + cfg.Env)

	// Start server
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}

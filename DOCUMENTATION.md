# E-Commerce Backend API - Complete Documentation

**Version:** 1.0.0  
**Author:** Jitender  
**Date:** February 2026  
**Tech Stack:** Go 1.24, Gin Framework, PostgreSQL 15, Docker

---

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [Architecture & Design Patterns](#2-architecture--design-patterns)
3. [Technology Stack](#3-technology-stack)
4. [Database Architecture](#4-database-architecture)
5. [Application Flow](#5-application-flow)
6. [Core Components](#6-core-components)
7. [API Endpoints](#7-api-endpoints)
8. [Security & Authentication](#8-security--authentication)
9. [Middleware](#9-middleware)
10. [Business Logic Flows](#10-business-logic-flows)
11. [Deployment](#11-deployment)
12. [Development Guide](#12-development-guide)

---

## 1. Project Overview

### 1.1 Description

A full-featured, production-ready e-commerce backend API built with modern Go practices. The system provides comprehensive functionality for managing products, shopping carts, orders, payments, and returns with robust authentication, authorization, and stock management.

### 1.2 Key Features

#### User Management

- ✅ User registration with email validation
- ✅ Secure authentication using JWT tokens
- ✅ Role-based access control (Customer, Admin)
- ✅ Profile management
- ✅ Password change functionality

#### Product Management

- ✅ Full CRUD operations for products
- ✅ Public product catalog with search and filtering
- ✅ Category-based organization
- ✅ Stock tracking and management
- ✅ Admin-only product management

#### Shopping Cart

- ✅ Add/Update/Remove items
- ✅ Real-time stock reservation system
- ✅ Cart validation before checkout
- ✅ Automatic stock release on cart expiry

#### Order Processing

- ✅ Order creation from cart
- ✅ Multiple payment methods (Credit Card, Debit Card, Cash on Delivery)
- ✅ Order status tracking (pending → processing → shipped → delivered → completed)
- ✅ Order cancellation
- ✅ Transaction management with database rollback

#### Payment System

- ✅ Mock payment gateway integration
- ✅ Payment status tracking
- ✅ Payment verification
- ✅ Automatic payment creation for card payments
- ✅ COD payment handling on delivery

#### Returns & Refunds

- ✅ Return request creation
- ✅ Return status management (requested → approved/rejected → completed)
- ✅ Refund processing
- ✅ Stock restoration on approved returns

#### Admin Dashboard

- ✅ User management
- ✅ Order management and analytics
- ✅ Product management
- ✅ Return processing
- ✅ Top products analysis
- ✅ Revenue tracking

---

## 2. Architecture & Design Patterns

### 2.1 Architectural Style

The application follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────┐
│                     HTTP Layer (Gin)                     │
│                    (Port: 8080)                          │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│                    Middleware Layer                      │
│  CORS | Auth | Logging | Recovery | Request ID          │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│                    Handlers Layer                        │
│  Request Validation | Response Formatting               │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│                    Service Layer                         │
│  Business Logic | Validation | Orchestration            │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│                  Repository Layer                        │
│  Data Access | SQL Queries | Transaction Management     │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│                PostgreSQL Database                       │
│                   (Port: 5433)                           │
└─────────────────────────────────────────────────────────┘
```

### 2.2 Design Patterns

#### 2.2.1 Repository Pattern

**Purpose:** Abstracts data access logic and provides a clean API for data operations.

**Implementation:**

```go
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    // ... more methods
}
```

**Benefits:**

- Testability: Easy to mock repositories for testing
- Flexibility: Can swap database implementations
- Single Responsibility: Each repository handles one entity

#### 2.2.2 Service Layer Pattern

**Purpose:** Encapsulates business logic and orchestrates operations across repositories.

**Implementation:**

```go
type OrderService interface {
    CreateOrder(ctx context.Context, userID uuid.UUID, req models.CreateOrderRequest) (*models.Order, error)
    GetOrder(ctx context.Context, orderID, userID uuid.UUID) (*models.Order, error)
    // ... more methods
}
```

**Responsibilities:**

- Business rule validation
- Transaction management
- Cross-entity operations
- Error handling and logging

#### 2.2.3 Dependency Injection

**Purpose:** Promotes loose coupling and testability.

**Implementation:**

```go
func InitRepositories(db *pgxpool.Pool, cfg *config.Config) *Repositories {
    // Initialize repositories
    userRepo := repository.NewUserRepository(db)
    productRepo := repository.NewProductRepository(db)

    // Initialize services with injected dependencies
    authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiry)
    productService := service.NewProductService(productRepo)

    // Initialize handlers with injected services
    authHandler := NewAuthHandler(authService)
    // ...
}
```

#### 2.2.4 Middleware Chain Pattern

**Purpose:** Provides cross-cutting concerns without polluting business logic.

```go
router.Use(middleware.GinCORSMiddleware(cfg.AllowedOrigins))
router.Use(middleware.NoCacheMiddleware())
router.Use(middleware.GinRecovery())
router.Use(middleware.GinLogging())
router.Use(middleware.GinRequestID())
```

### 2.3 Project Structure

```
ecommerce-backend-go/
│
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
│
├── internal/
│   ├── config/
│   │   └── config.go               # Configuration management
│   │
│   ├── handlers/                   # HTTP request handlers
│   │   ├── init.go                 # Dependency injection setup
│   │   ├── auth.go                 # Authentication handlers
│   │   ├── cart.go                 # Shopping cart handlers
│   │   ├── health.go               # Health check handlers
│   │   ├── orders.go               # Order management handlers
│   │   ├── payments.go             # Payment handlers
│   │   ├── products.go             # Product handlers
│   │   └── returns.go              # Return handlers
│   │
│   ├── middleware/                 # HTTP middleware
│   │   ├── auth.go                 # Legacy auth middleware
│   │   ├── cors.go                 # CORS configuration
│   │   ├── gin_auth.go             # Gin authentication
│   │   ├── gin_common.go           # Common Gin middleware
│   │   ├── logging.go              # Request logging
│   │   ├── rate_limit.go           # Rate limiting
│   │   └── recovery.go             # Panic recovery
│   │
│   ├── models/                     # Domain models
│   │   ├── admin.go                # Admin-specific models
│   │   ├── cart.go                 # Cart models
│   │   ├── order.go                # Order models
│   │   ├── payment.go              # Payment models
│   │   ├── product.go              # Product models
│   │   ├── return.go               # Return models
│   │   └── user.go                 # User models
│   │
│   ├── repository/                 # Data access layer
│   │   ├── database.go             # Database utilities
│   │   ├── cart_repo.go            # Cart data access
│   │   ├── order_repo.go           # Order data access
│   │   ├── payment_repo.go         # Payment data access
│   │   ├── product_repo.go         # Product data access
│   │   ├── return_repo.go          # Return data access
│   │   └── user_repo.go            # User data access
│   │
│   └── service/                    # Business logic layer
│       ├── auth_service.go         # Authentication service
│       ├── cart_service.go         # Cart service
│       ├── order_service.go        # Order service
│       ├── payment_service.go      # Payment service
│       ├── product_service.go      # Product service
│       └── return_service.go       # Return service
│
├── pkg/                            # Public packages
│   ├── database/
│   │   └── postgres.go             # PostgreSQL connection
│   │
│   └── utils/                      # Utility functions
│       ├── gin_response.go         # Gin response helpers
│       ├── jwt.go                  # JWT utilities
│       ├── password.go             # Password hashing
│       ├── response.go             # Response formatting
│       └── validators.go           # Input validation
│
├── migrations/
│   └── 001_init.sql                # Database schema
│
├── docker-compose.yml              # Docker orchestration
├── Dockerfile                      # Container definition
├── go.mod                          # Go dependencies
├── go.sum                          # Dependency checksums
├── Makefile                        # Build automation
├── openapi.yaml                    # API specification
├── postman_collection.json         # API testing collection
├── README.md                       # Project readme
└── SETUP_AND_RUN.md               # Setup instructions
```

---

## 3. Technology Stack

### 3.1 Core Technologies

| Component        | Technology     | Version | Purpose                                  |
| ---------------- | -------------- | ------- | ---------------------------------------- |
| Language         | Go             | 1.24    | Backend development                      |
| Web Framework    | Gin            | 1.10.0  | HTTP routing and middleware              |
| Database         | PostgreSQL     | 15      | Data persistence                         |
| Database Driver  | pgx            | 5.8.0   | PostgreSQL driver and connection pooling |
| Authentication   | JWT            | 5.3.0   | Token-based auth                         |
| Validation       | validator      | 10.30.1 | Input validation                         |
| Containerization | Docker         | Latest  | Application containerization             |
| Orchestration    | Docker Compose | 3.8     | Multi-container setup                    |

### 3.2 Key Dependencies

```go
require (
    github.com/gin-gonic/gin v1.10.0              // Web framework
    github.com/go-playground/validator/v10 v10.30.1  // Validation
    github.com/golang-jwt/jwt/v5 v5.3.0           // JWT tokens
    github.com/google/uuid v1.6.0                 // UUID generation
    github.com/jackc/pgx/v5 v5.8.0                // PostgreSQL driver
    github.com/joho/godotenv v1.5.1               // Environment variables
    golang.org/x/crypto v0.47.0                   // Password hashing
)
```

### 3.3 Infrastructure

#### Development Environment

- **OS:** Windows (Docker Desktop)
- **Port Mapping:**
  - Backend API: `8080`
  - PostgreSQL: `5433:5432`

#### Production Considerations

- Connection pooling configured
- Health checks enabled
- Graceful shutdown support
- Environment-based configuration

---

## 4. Database Architecture

### 4.1 Entity Relationship Diagram

```
┌──────────────┐           ┌──────────────┐
│    users     │           │   products   │
├──────────────┤           ├──────────────┤
│ id (PK)      │           │ id (PK)      │
│ email        │           │ sku (UNIQUE) │
│ password_hash│           │ name         │
│ first_name   │           │ description  │
│ last_name    │           │ price        │
│ role         │           │ stock_qty    │
│ created_at   │           │ category     │
│ updated_at   │           │ image_url    │
└──────────────┘           │ created_at   │
       │                   │ updated_at   │
       │                   └──────────────┘
       │                          │
       │ 1                        │
       │                          │ *
       │ *                   ┌────┴─────────┐
       ├─────────────────────┤  cart_items  │
       │                     ├──────────────┤
       │                     │ id (PK)      │
       │                     │ cart_id (FK) │
       │                     │ product_id   │
       │                     │ quantity     │
       │                     │ created_at   │
       │                     └──────────────┘
       │                          │
       │ 1                        │ *
       │                     ┌────┴──────┐
       ├─────────────────────┤   carts   │
       │                     ├───────────┤
       │                     │ id (PK)   │
       │                     │ user_id   │
       │                     │ created_at│
       │                     │ updated_at│
       │                     └───────────┘
       │
       │ 1
       │
       │ *
  ┌────┴──────────┐
  │    orders     │
  ├───────────────┤
  │ id (PK)       │
  │ user_id (FK)  │
  │ order_number  │
  │ total_amount  │
  │ status        │
  │ payment_method│
  │ ship_address  │
  │ bill_address  │
  │ created_at    │
  │ updated_at    │
  └───────────────┘
       │
       │ 1
       │
       ├──────────────────────┐
       │ *                    │ 1
  ┌────┴──────────┐      ┌────┴──────────┐
  │ order_items   │      │   payments    │
  ├───────────────┤      ├───────────────┤
  │ id (PK)       │      │ id (PK)       │
  │ order_id (FK) │      │ order_id (FK) │
  │ product_id    │      │ amount        │
  │ quantity      │      │ status        │
  │ price_at_time │      │ payment_method│
  │ created_at    │      │ transaction_id│
  └───────────────┘      │ details (JSON)│
                         │ created_at    │
       │                 │ updated_at    │
       │ 1               └───────────────┘
       │
       │ *
  ┌────┴──────────┐
  │    returns    │
  ├───────────────┤
  │ id (PK)       │
  │ order_id (FK) │
  │ user_id (FK)  │
  │ reason        │
  │ status        │
  │ refund_amount │
  │ created_at    │
  │ updated_at    │
  └───────────────┘


  ┌──────────────────────┐
  │ stock_reservations   │
  ├──────────────────────┤
  │ id (PK)              │
  │ product_id (FK)      │
  │ cart_id (FK)         │
  │ quantity             │
  │ expires_at           │
  │ created_at           │
  └──────────────────────┘
```

### 4.2 Table Descriptions

#### users

**Purpose:** Store user account information

- **Primary Key:** `id` (UUID)
- **Unique Constraints:** `email`
- **Check Constraints:** `role IN ('customer', 'admin')`
- **Indexes:** Primary key index

#### products

**Purpose:** Store product catalog

- **Primary Key:** `id` (UUID)
- **Unique Constraints:** `sku`
- **Check Constraints:** `price >= 0`, `stock_quantity >= 0`
- **Indexes:** `sku`, `category`

#### carts

**Purpose:** Store user shopping carts

- **Primary Key:** `id` (UUID)
- **Foreign Keys:** `user_id` → `users(id)` ON DELETE CASCADE
- **Unique Constraints:** `user_id` (one cart per user)

#### cart_items

**Purpose:** Store items in shopping carts

- **Primary Key:** `id` (UUID)
- **Foreign Keys:**
  - `cart_id` → `carts(id)` ON DELETE CASCADE
  - `product_id` → `products(id)` ON DELETE CASCADE
- **Unique Constraints:** `(cart_id, product_id)`
- **Check Constraints:** `quantity > 0`

#### orders

**Purpose:** Store customer orders

- **Primary Key:** `id` (UUID)
- **Foreign Keys:** `user_id` → `users(id)`
- **Unique Constraints:** `order_number`
- **Check Constraints:**
  - `total_amount >= 0`
  - `status IN ('pending', 'processing', 'shipped', 'delivered', 'completed', 'cancelled', 'refunded', 'return_requested')`
  - `payment_method IN ('cc', 'dc', 'cod')`
- **Indexes:** `user_id`, `status`

#### order_items

**Purpose:** Store items in orders

- **Primary Key:** `id` (UUID)
- **Foreign Keys:**
  - `order_id` → `orders(id)` ON DELETE CASCADE
  - `product_id` → `products(id)`
- **Check Constraints:** `quantity > 0`, `price_at_time >= 0`
- **Indexes:** `order_id`

#### payments

**Purpose:** Track payment transactions

- **Primary Key:** `id` (UUID)
- **Foreign Keys:** `order_id` → `orders(id)` ON DELETE CASCADE
- **Check Constraints:**
  - `amount >= 0`
  - `status IN ('pending', 'processing', 'completed', 'failed', 'refunded')`
- **Indexes:** `order_id`

#### returns

**Purpose:** Handle product returns and refunds

- **Primary Key:** `id` (UUID)
- **Foreign Keys:**
  - `order_id` → `orders(id)`
  - `user_id` → `users(id)`
- **Check Constraints:** `status IN ('requested', 'approved', 'rejected', 'completed')`

#### stock_reservations

**Purpose:** Prevent overselling by reserving stock temporarily

- **Primary Key:** `id` (UUID)
- **Foreign Keys:**
  - `product_id` → `products(id)` ON DELETE CASCADE
  - `cart_id` → `carts(id)` ON DELETE CASCADE
- **Unique Constraints:** `(product_id, cart_id)`
- **Check Constraints:** `quantity > 0`
- **Indexes:** `expires_at`
- **Expiry:** Automatically expires based on `expires_at` timestamp

### 4.3 Database Triggers

```sql
-- Auto-update timestamp triggers
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_carts_updated_at
    BEFORE UPDATE ON carts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payments_updated_at
    BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_returns_updated_at
    BEFORE UPDATE ON returns
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

---

## 5. Application Flow

### 5.1 Application Startup Sequence

```
┌─────────────────────────────────────────────────────────┐
│ 1. Load Configuration                                   │
│    - Environment variables                              │
│    - .env file                                          │
│    - Default values                                     │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 2. Initialize Database Connection                       │
│    - Build connection string                            │
│    - Configure connection pool                          │
│    - Test connection with ping                          │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 3. Initialize Dependencies (Dependency Injection)       │
│    - Create repositories (data access layer)            │
│    - Create services (business logic layer)             │
│    - Create handlers (HTTP layer)                       │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 4. Setup HTTP Router (Gin)                             │
│    - Apply global middleware                            │
│    - Register public routes                             │
│    - Register protected routes (with auth)              │
│    - Register admin routes (with auth + admin check)    │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 5. Start HTTP Server                                    │
│    - Listen on configured port (default: 8080)          │
│    - Ready to accept requests                           │
└─────────────────────────────────────────────────────────┘
```

### 5.2 Request Lifecycle

```
┌──────────────────────────────────────────────────────────┐
│ CLIENT REQUEST                                           │
└────────────────────┬─────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ HTTP Layer: Gin Router                                  │
│  - Route matching                                       │
│  - HTTP method validation                               │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ Middleware Chain (Executed in Order)                    │
│  1. CORS            - Set CORS headers                  │
│  2. RequestID       - Generate unique request ID        │
│  3. Logging         - Log request details               │
│  4. Recovery        - Catch panics                      │
│  5. NoCache         - Set no-cache headers              │
│  6. Auth (optional) - Validate JWT token                │
│  7. Admin (optional)- Verify admin role                 │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ Handler Layer                                           │
│  - Parse request body (JSON)                            │
│  - Validate input using validator                       │
│  - Extract user context (ID, role)                      │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ Service Layer                                           │
│  - Execute business logic                               │
│  - Validate business rules                              │
│  - Coordinate multiple repositories                     │
│  - Handle transactions                                  │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ Repository Layer                                        │
│  - Execute SQL queries                                  │
│  - Handle database connections                          │
│  - Manage transactions                                  │
│  - Map database rows to Go structs                      │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ PostgreSQL Database                                     │
│  - Execute query                                        │
│  - Apply constraints                                    │
│  - Return results                                       │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ Response Flow (Reverse Direction)                       │
│  Repository → Service → Handler → Middleware → Client   │
│  - Format response                                      │
│  - Set HTTP status code                                 │
│  - Send JSON response                                   │
└──────────────────────────────────────────────────────────┘
```

### 5.3 Authentication Flow

```
┌──────────────────────────────────────────────────────────┐
│ 1. USER LOGIN REQUEST                                    │
│    POST /api/v1/auth/login                              │
│    Body: {email, password}                              │
└────────────────────┬─────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 2. AuthHandler.Login                                    │
│    - Validate request format                            │
│    - Call AuthService.Login()                           │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 3. AuthService.Login                                    │
│    - Get user by email from UserRepository              │
│    - Verify password hash (bcrypt)                      │
│    - Generate JWT token                                 │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 4. JWT Token Generation                                 │
│    Claims:                                              │
│    - user_id: UUID                                      │
│    - email: string                                      │
│    - role: "customer" | "admin"                         │
│    - exp: expiration timestamp                          │
│    - iat: issued at timestamp                           │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 5. RESPONSE TO CLIENT                                   │
│    {                                                    │
│      "success": true,                                   │
│      "message": "Login successful",                     │
│      "data": {                                          │
│        "user": {...},                                   │
│        "access_token": "eyJhbGc..."                     │
│      }                                                  │
│    }                                                    │
└─────────────────────────────────────────────────────────┘


┌──────────────────────────────────────────────────────────┐
│ 6. SUBSEQUENT AUTHENTICATED REQUESTS                     │
│    Header: Authorization: Bearer <token>                │
└────────────────────┬─────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 7. GinAuthMiddleware                                    │
│    - Extract Bearer token from Authorization header     │
│    - Validate JWT signature                             │
│    - Check expiration                                   │
│    - Extract claims (user_id, role)                     │
│    - Set in Gin context                                 │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 8. Handler Access to User Context                       │
│    userID := c.Get("userID")                            │
│    role := c.Get("userRole")                            │
└─────────────────────────────────────────────────────────┘
```

---

## 6. Core Components

### 6.1 Models Layer

**Location:** `internal/models/`

Models define the domain entities and data structures used throughout the application.

#### User Model

```go
type User struct {
    ID           uuid.UUID `json:"id"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`          // Hidden from JSON
    FirstName    string    `json:"first_name"`
    LastName     string    `json:"last_name"`
    Role         string    `json:"role"`       // "customer" or "admin"
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

#### Product Model

```go
type Product struct {
    ID          uuid.UUID `json:"id"`
    SKU         string    `json:"sku"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    Stock       int       `json:"stock"`       // Available stock
    Category    string    `json:"category"`
    ImageURL    string    `json:"image_url"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### Cart Model

```go
type Cart struct {
    ID        uuid.UUID  `json:"id"`
    UserID    uuid.UUID  `json:"user_id"`
    Items     []CartItem `json:"items"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}

type CartItem struct {
    ID        uuid.UUID `json:"id"`
    CartID    uuid.UUID `json:"cart_id"`
    ProductID uuid.UUID `json:"product_id"`
    Product   Product   `json:"product"`
    Quantity  int       `json:"quantity"`
    CreatedAt time.Time `json:"created_at"`
}
```

#### Order Model

```go
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

type OrderStatus string
const (
    OrderPending         OrderStatus = "pending"
    OrderProcessing      OrderStatus = "processing"
    OrderShipped         OrderStatus = "shipped"
    OrderDelivered       OrderStatus = "delivered"
    OrderCompleted       OrderStatus = "completed"
    OrderCancelled       OrderStatus = "cancelled"
    OrderRefunded        OrderStatus = "refunded"
    OrderReturnRequested OrderStatus = "return_requested"
)
```

### 6.2 Repository Layer

**Location:** `internal/repository/`

Repositories handle all database operations and provide a clean API for data access.

#### Key Responsibilities:

1. Execute SQL queries
2. Map database rows to Go structs
3. Handle database errors
4. Manage transactions
5. Implement connection pooling

#### Example: UserRepository

```go
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    GetAll(ctx context.Context, page, limit, rangeDays int) ([]models.User, int, error)
    Update(ctx context.Context, user *models.User) error
    UpdateRole(ctx context.Context, id uuid.UUID, role string) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

#### Stock Reservation System

The ProductRepository implements a sophisticated stock reservation system:

```go
// Reserve stock when adding to cart
func ReserveStock(ctx context.Context, productID, cartID uuid.UUID,
                  quantity int, expiresAt int64) error

// Calculate available stock excluding reservations
func GetAvailableStock(ctx context.Context, productID uuid.UUID) (int, error)

// Calculate available stock excluding specific cart's reservations
func GetAvailableStockExcludingCart(ctx context.Context,
                                   productID, cartID uuid.UUID) (int, error)
```

### 6.3 Service Layer

**Location:** `internal/service/`

Services contain business logic and orchestrate operations across multiple repositories.

#### Key Responsibilities:

1. Implement business rules
2. Validate business constraints
3. Coordinate multiple repositories
4. Handle complex transactions
5. Log business events

#### Example: OrderService

```go
type OrderService interface {
    CreateOrder(ctx context.Context, userID uuid.UUID,
                req models.CreateOrderRequest) (*models.Order, error)
    GetOrder(ctx context.Context, orderID, userID uuid.UUID) (*models.Order, error)
    GetUserOrders(ctx context.Context, userID uuid.UUID,
                  page, limit int) ([]models.Order, int, error)
    UpdateOrderStatus(ctx context.Context, orderID uuid.UUID,
                      status models.OrderStatus) error
    CancelOrder(ctx context.Context, orderID, userID uuid.UUID) error
}
```

**CreateOrder Business Logic:**

1. Validate cart is not empty
2. Validate cart stock availability
3. Begin database transaction
4. Calculate total amount
5. Deduct stock from inventory
6. Create order with items
7. Clear user's cart
8. Commit transaction
9. Create payment (if card payment)
10. Return created order

### 6.4 Handler Layer

**Location:** `internal/handlers/`

Handlers process HTTP requests and return responses.

#### Key Responsibilities:

1. Parse request body
2. Validate input
3. Extract user context
4. Call service methods
5. Format responses
6. Handle errors

#### Response Format

```go
// Success Response
{
    "success": true,
    "message": "Operation successful",
    "data": {...}
}

// Error Response
{
    "success": false,
    "message": "Operation failed",
    "error": "Error details"
}

// Validation Error Response
{
    "success": false,
    "message": "Validation failed",
    "errors": [
        {"field": "email", "message": "Invalid email format"}
    ]
}
```

---

## 7. API Endpoints

### 7.1 Public Endpoints

#### Health Checks

```
GET  /health              - Basic health check
GET  /ready               - Readiness check
GET  /metrics             - System metrics
GET  /api/v1/health       - API version health check
GET  /api/v1/ready        - API version readiness check
GET  /api/v1/metrics      - API version metrics
```

#### Authentication

```
POST /api/v1/auth/register   - Register new user
POST /api/v1/auth/login      - Login user
POST /api/v1/auth/refresh    - Refresh access token
```

#### Products (Public)

```
GET  /api/v1/products         - List products (with pagination, search, filter)
GET  /api/v1/products/:id     - Get single product
```

### 7.2 Protected Endpoints (Require Authentication)

#### User Profile

```
GET  /api/v1/users/profile           - Get current user profile
PUT  /api/v1/users/profile           - Update profile
PUT  /api/v1/users/change-password   - Change password
```

#### Shopping Cart

```
GET    /api/v1/cart                - Get user's cart
GET    /api/v1/cart/validate       - Validate cart before checkout
POST   /api/v1/cart/items          - Add item to cart
PUT    /api/v1/cart/items/:itemId  - Update cart item quantity
DELETE /api/v1/cart/items/:itemId  - Remove item from cart
DELETE /api/v1/cart                - Clear entire cart
```

#### Orders

```
POST   /api/v1/orders              - Create order from cart
GET    /api/v1/orders              - Get user's orders (paginated)
GET    /api/v1/orders/:id          - Get specific order
GET    /api/v1/orders/:id/payment  - Get payment details for order
PUT    /api/v1/orders/:id/cancel   - Cancel order
```

#### Payments

```
POST   /api/v1/payments              - Create payment
POST   /api/v1/payments/:id/verify   - Verify payment
```

#### Returns

```
POST   /api/v1/returns     - Create return request
GET    /api/v1/returns     - Get user's returns
GET    /api/v1/returns/:id - Get specific return details
```

### 7.3 Admin Endpoints (Require Admin Role)

#### Product Management

```
POST   /api/v1/admin/products      - Create product
GET    /api/v1/admin/products      - Get all products (admin view)
PUT    /api/v1/admin/products/:id  - Update product
DELETE /api/v1/admin/products/:id  - Delete product
GET    /api/v1/admin/products/top  - Get top selling products
```

#### Order Management

```
GET  /api/v1/admin/orders            - Get all orders (paginated, filtered)
GET  /api/v1/admin/orders/recent     - Get recent orders
GET  /api/v1/admin/orders/:id        - Get order details (admin view)
PUT  /api/v1/admin/orders/:id/status - Update order status
GET  /api/v1/admin/analytics         - Get sales analytics
```

#### User Management

```
GET  /api/v1/admin/users         - Get all users
PUT  /api/v1/admin/users/:id/role - Update user role
```

#### Return Management

```
GET  /api/v1/admin/returns                   - Get all returns
POST /api/v1/admin/returns/:returnId/process - Process return (approve/reject)
```

### 7.4 Query Parameters

#### Pagination

```
?page=1&limit=10
```

#### Filtering

```
?category=Electronics
?search=laptop
?status=pending
?range_days=30
```

---

## 8. Security & Authentication

### 8.1 Authentication Mechanism

The application uses **JWT (JSON Web Tokens)** for authentication.

#### Token Structure

```
Header: {
    "alg": "HS256",
    "typ": "JWT"
}

Payload: {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "role": "customer",
    "exp": 1709856000,
    "iat": 1709769600
}

Signature: HMACSHA256(
    base64UrlEncode(header) + "." + base64UrlEncode(payload),
    JWT_SECRET
)
```

#### Token Configuration

- **Algorithm:** HS256 (HMAC with SHA-256)
- **Expiry:** 24 hours (configurable via `JWT_EXPIRY_HOURS`)
- **Secret:** Configured via environment variable `JWT_SECRET`

### 8.2 Password Security

#### Hashing

- **Algorithm:** bcrypt
- **Cost Factor:** bcrypt.DefaultCost (10)
- **Salt:** Automatically generated per password

```go
// Hash password during registration
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(password),
    bcrypt.DefaultCost
)

// Verify password during login
err := bcrypt.CompareHashAndPassword(
    []byte(hash),
    []byte(password)
)
```

### 8.3 Authorization

#### Role-Based Access Control (RBAC)

**Roles:**

1. **customer** (default) - Regular users
2. **admin** - Administrators

**Middleware Chain:**

```go
// Protected routes (requires authentication)
protected := api.Group("")
protected.Use(middleware.GinAuthMiddleware(authService))

// Admin routes (requires authentication + admin role)
admin := api.Group("/admin")
admin.Use(middleware.GinAuthMiddleware(authService))
admin.Use(middleware.GinAdminMiddleware())
```

#### Admin-Only Operations

- Create/Update/Delete products
- View all orders
- Update order status
- Process returns
- View analytics
- Manage users
- Change user roles

### 8.4 Security Best Practices Implemented

1. **Password Never Stored in Plain Text**
   - Always hashed with bcrypt
   - Password hash never returned in API responses

2. **SQL Injection Prevention**
   - Parameterized queries using pgx
   - No string concatenation for SQL queries

3. **CORS Configuration**
   - Configurable allowed origins
   - Credentials support
   - Preflight request handling

4. **Input Validation**
   - Request body validation using validator library
   - UUID validation
   - Email format validation
   - Required field validation

5. **Error Handling**
   - No sensitive information in error messages
   - Generic error responses to clients
   - Detailed logging server-side

6. **Rate Limiting**
   - Middleware support included
   - Configurable per endpoint

---

## 9. Middleware

### 9.1 CORS Middleware

**Purpose:** Handle Cross-Origin Resource Sharing

**Configuration:**

```go
router.Use(middleware.GinCORSMiddleware(cfg.AllowedOrigins))
```

**Functionality:**

- Validates origin against allowed list
- Sets appropriate CORS headers
- Handles preflight OPTIONS requests
- Configurable allowed origins via environment

### 9.2 Authentication Middleware

**Purpose:** Validate JWT tokens and extract user information

**File:** `internal/middleware/gin_auth.go`

**Flow:**

1. Extract Authorization header
2. Validate Bearer token format
3. Parse and verify JWT signature
4. Check token expiration
5. Extract claims (user_id, role)
6. Store in Gin context for handlers

**Usage:**

```go
protected.Use(middleware.GinAuthMiddleware(authService))
```

### 9.3 Admin Middleware

**Purpose:** Verify user has admin role

**File:** `internal/middleware/gin_auth.go`

**Flow:**

1. Check if user role exists in context
2. Verify role == "admin"
3. Return 403 Forbidden if not admin

**Usage:**

```go
admin.Use(middleware.GinAdminMiddleware())
```

### 9.4 Logging Middleware

**Purpose:** Log all incoming requests

**Functionality:**

- Log HTTP method and path
- Log status code
- Log response time
- Log request ID

### 9.5 Recovery Middleware

**Purpose:** Recover from panics and return 500 error

**Functionality:**

- Catch panics in handlers
- Log panic details
- Return generic 500 error to client
- Prevent server crash

### 9.6 Request ID Middleware

**Purpose:** Generate unique ID for each request

**Functionality:**

- Generate UUID for each request
- Store in Gin context
- Include in logs
- Return in response headers

### 9.7 No-Cache Middleware

**Purpose:** Prevent caching of API responses

**Headers Set:**

```
Cache-Control: no-cache, no-store, must-revalidate
Pragma: no-cache
Expires: 0
```

---

## 10. Business Logic Flows

### 10.1 User Registration Flow

```
┌─────────────────────────────────────────────────────────┐
│ 1. Client sends POST /api/v1/auth/register              │
│    {email, password, first_name, last_name}             │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 2. AuthHandler validates input                          │
│    - Email format                                       │
│    - Password length (min 6 chars)                      │
│    - Required fields present                            │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 3. AuthService.Register                                 │
│    a) Check if email already exists                     │
│    b) Hash password using bcrypt                        │
│    c) Create user with role="customer"                  │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 4. UserRepository.Create                                │
│    - Insert into database                               │
│    - Return generated UUID and timestamps               │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 5. Generate JWT token                                   │
│    - Claims: user_id, email, role                       │
│    - Expiry: 24 hours                                   │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 6. Return response                                      │
│    {user: {...}, access_token: "..."}                   │
└─────────────────────────────────────────────────────────┘
```

### 10.2 Add to Cart Flow

```
┌─────────────────────────────────────────────────────────┐
│ 1. Client sends POST /api/v1/cart/items                │
│    Header: Authorization: Bearer <token>                │
│    Body: {product_id, quantity}                         │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 2. Authentication Middleware                            │
│    - Validate JWT token                                 │
│    - Extract user_id                                    │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 3. CartHandler validates input                          │
│    - Valid UUID for product_id                          │
│    - Quantity > 0                                       │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 4. CartService.AddToCart                                │
│    a) Get or create user's cart                         │
│    b) Verify product exists                             │
│    c) Check available stock                             │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 5. ProductService.CheckStock                            │
│    - Calculate: actual_stock - reserved_stock           │
│    - Compare with requested quantity                    │
└────────────────────┬────────────────────────────────────┘
                     ↓ (if sufficient)
┌─────────────────────────────────────────────────────────┐
│ 6. ProductService.ReserveStock                          │
│    - Insert into stock_reservations table               │
│    - Set expires_at = NOW() + TTL (10 minutes)          │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 7. CartRepository.AddItem                               │
│    - Insert into cart_items table                       │
│    - If product already in cart, update quantity        │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 8. Return updated cart with all items                   │
└─────────────────────────────────────────────────────────┘
```

### 10.3 Create Order Flow (Most Complex)

```
┌─────────────────────────────────────────────────────────┐
│ 1. Client sends POST /api/v1/orders                     │
│    {shipping_address, billing_address, payment_method}  │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 2. Authentication & Validation                          │
│    - Validate JWT token                                 │
│    - Validate payment_method: cc|dc|cod                 │
│    - Validate address fields                            │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 3. OrderService.CreateOrder                             │
│    a) Get user's cart with items                        │
│    b) Check cart is not empty                           │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 4. CartService.ValidateCart                             │
│    - For each item:                                     │
│      * Check product still exists                       │
│      * Verify stock availability                        │
│      * Check price hasn't changed (optional warning)    │
└────────────────────┬────────────────────────────────────┘
                     ↓ (if valid)
┌─────────────────────────────────────────────────────────┐
│ 5. BEGIN DATABASE TRANSACTION                           │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 6. Calculate Total & Prepare Order Items                │
│    - For each cart item:                                │
│      * Calculate: price * quantity                      │
│      * Create OrderItem with price_at_time              │
│      * Add to total_amount                              │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 7. Deduct Stock from Inventory                          │
│    - For each item:                                     │
│      UPDATE products                                    │
│      SET stock_quantity = stock_quantity - quantity     │
│      WHERE id = product_id                              │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 8. Create Order Record                                  │
│    - Generate unique order_number                       │
│    - Set status = "pending"                             │
│    - Insert into orders table                           │
│    - Insert order_items                                 │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 9. Clear User's Cart                                    │
│    - Delete from cart_items                             │
│    - Delete stock reservations                          │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 10. COMMIT TRANSACTION                                  │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 11. Create Payment (if CC/DC)                           │
│     - PaymentService.CreatePaymentForOrder              │
│     - Insert into payments table                        │
│     - Status = "completed" (mock)                       │
│     - Update order status = "processing"                │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 12. Return Created Order                                │
│     - Include all order items                           │
│     - Include addresses                                 │
│     - Include order_number                              │
└─────────────────────────────────────────────────────────┘

Note: If COD, payment is created when status changes to "delivered"
```

### 10.4 Update Order Status Flow

```
┌─────────────────────────────────────────────────────────┐
│ 1. Admin sends PUT /api/v1/admin/orders/:id/status      │
│    Body: {status: "shipped"}                            │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 2. Authentication & Authorization                       │
│    - Validate JWT token                                 │
│    - Verify role = "admin"                              │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 3. OrderService.UpdateOrderStatus                       │
│    a) Get current order                                 │
│    b) Validate status transition                        │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 4. Status Transition Validation                         │
│    Valid transitions:                                   │
│    pending → processing, cancelled                      │
│    processing → shipped, delivered, completed, cancelled│
│    shipped → delivered, completed                       │
│    delivered → completed                                │
│    completed → (none)                                   │
│    cancelled → (none)                                   │
└────────────────────┬────────────────────────────────────┘
                     ↓ (if valid)
┌─────────────────────────────────────────────────────────┐
│ 5. Update Order Status in Database                      │
│    UPDATE orders SET status = ?, updated_at = NOW()     │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 6. Special Handling for "delivered" + COD               │
│    IF status = "delivered" AND payment_method = "cod"   │
│    THEN create payment record                           │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 7. Return Success Response                              │
└─────────────────────────────────────────────────────────┘
```

### 10.5 Return Request Flow

```
┌─────────────────────────────────────────────────────────┐
│ 1. Customer sends POST /api/v1/returns                  │
│    Body: {order_id, reason}                             │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 2. ReturnService.CreateReturn                           │
│    a) Verify order exists                               │
│    b) Verify order belongs to user                      │
│    c) Check order is eligible for return                │
│       (status = delivered or completed)                 │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 3. Create Return Record                                 │
│    - status = "requested"                               │
│    - Insert into returns table                          │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 4. Update Order Status                                  │
│    - Set order.status = "return_requested"              │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 5. Admin Reviews Return                                 │
│    POST /api/v1/admin/returns/:id/process               │
│    Body: {status: "approved", refund_amount: 150.00}    │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 6. ReturnService.ProcessReturn                          │
│    IF status = "approved":                              │
│      a) Update return status and refund_amount          │
│      b) Process refund through PaymentService           │
│      c) Update payment status = "refunded"              │
│      d) Restore stock to inventory                      │
│      e) Update order status = "refunded"                │
│    IF status = "rejected":                              │
│      a) Update return status                            │
│      b) Restore order to previous status                │
└─────────────────────────────────────────────────────────┘
```

### 10.6 Stock Reservation System

**Purpose:** Prevent overselling when multiple users add the same product to cart simultaneously.

```
┌─────────────────────────────────────────────────────────┐
│ Product Stock: 100 units                                │
│ User A adds 80 to cart                                  │
│ User B adds 30 to cart (simultaneously)                 │
└─────────────────────────────────────────────────────────┘

WITHOUT RESERVATION SYSTEM:
┌─────────────────────────────────────────────────────────┐
│ User A checks stock: 100 available ✓                    │
│ User B checks stock: 100 available ✓                    │
│ Both add to cart successfully!                          │
│ Total reserved: 110 (OVERSOLD by 10!)                   │
└─────────────────────────────────────────────────────────┘

WITH RESERVATION SYSTEM:
┌─────────────────────────────────────────────────────────┐
│ User A:                                                 │
│   1. Check: 100 - 0 reserved = 100 available ✓          │
│   2. Reserve 80 units (expires in 10 min)               │
│   3. Add to cart SUCCESS                                │
│                                                         │
│ User B (milliseconds later):                            │
│   1. Check: 100 - 80 reserved = 20 available            │
│   2. Request 30 > 20 available ✗                        │
│   3. Error: "Insufficient stock"                        │
└─────────────────────────────────────────────────────────┘

RESERVATION LIFECYCLE:
┌─────────────────────────────────────────────────────────┐
│ Add to Cart → Reserve Stock (TTL: 10 min)              │
│ Update Cart → Adjust Reservation                        │
│ Remove Item → Release Reservation                       │
│ Create Order → Deduct Stock, Delete Reservation         │
│ Expiry      → Auto-release after 10 min                 │
└─────────────────────────────────────────────────────────┘
```

**Database Query for Available Stock:**

```sql
SELECT
    p.stock_quantity - COALESCE(SUM(sr.quantity), 0) as available_stock
FROM products p
LEFT JOIN stock_reservations sr
    ON p.id = sr.product_id
    AND sr.expires_at > NOW()
WHERE p.id = $1
GROUP BY p.id, p.stock_quantity
```

---

## 11. Deployment

### 11.1 Docker Compose Deployment (Recommended)

**File:** `docker-compose.yml`

```yaml
version: "3.8"

services:
  postgres:
    image: postgres:15-alpine
    container_name: ecommerce-db
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: Jitender@123
      POSTGRES_DB: ecommerce_db
    ports:
      - "5433:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  backend:
    build: .
    container_name: ecommerce-backend
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: Jitender@123
      DB_NAME: ecommerce_db
      JWT_SECRET: your-secret-key
    ports:
      - "8080:8080"

volumes:
  postgres_data:
```

**Startup Commands:**

```bash
# Start services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### 11.2 Dockerfile

```dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
```

### 11.3 Environment Variables

**Required:**

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name (default: ecommerce_db)
- `JWT_SECRET` - Secret key for JWT signing

**Optional:**

- `PORT` - Server port (default: 8080)
- `ENV` - Environment (development/production)
- `JWT_EXPIRY_HOURS` - Token expiry (default: 24)
- `ALLOWED_ORIGINS` - CORS allowed origins (comma-separated)
- `STOCK_RESERVATION_TTL_MINUTES` - Stock reservation timeout (default: 10)
- `DB_SSLMODE` - PostgreSQL SSL mode (default: disable)

### 11.4 Database Migration

**Automatic Migration on Container Start:**
The PostgreSQL container automatically runs migration scripts from `./migrations/` directory on first startup.

**Manual Migration:**

```bash
# Connect to database
psql -h localhost -p 5433 -U postgres -d ecommerce_db

# Run migration script
\i migrations/001_init.sql
```

### 11.5 Health Checks

**Application Health:**

```bash
curl http://localhost:8080/health
```

**Response:**

```json
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2026-02-07T10:30:00Z"
}
```

**Readiness Check:**

```bash
curl http://localhost:8080/ready
```

---

## 12. Development Guide

### 12.1 Prerequisites

- Go 1.24+
- PostgreSQL 15+
- Docker & Docker Compose
- Git

### 12.2 Local Setup

```bash
# Clone repository
git clone <repository-url>
cd ecommerce-backend-go

# Install dependencies
go mod download

# Create .env file
cat > .env << EOF
PORT=8080
ENV=development
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=Jitender@123
DB_NAME=ecommerce_db
JWT_SECRET=dev-secret-key-change-in-production
JWT_EXPIRY_HOURS=24
ALLOWED_ORIGINS=http://localhost:3000
STOCK_RESERVATION_TTL_MINUTES=10
EOF

# Start PostgreSQL with Docker Compose
docker-compose up postgres -d

# Wait for database to be ready
sleep 10

# Run application
go run cmd/server/main.go
```

### 12.3 Database Setup

```bash
# Start PostgreSQL
docker-compose up postgres -d

# Run migrations
psql -h localhost -p 5433 -U postgres -d ecommerce_db -f migrations/001_init.sql

# Verify tables
psql -h localhost -p 5433 -U postgres -d ecommerce_db -c "\dt"
```

### 12.4 Testing

**Test Default Admin Login:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

**Test User Registration:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

**Test Product Listing:**

```bash
curl http://localhost:8080/api/v1/products
```

### 12.5 Makefile Targets

```makefile
.PHONY: run build test docker-up docker-down migrate

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

test:
	go test -v ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate:
	psql -h localhost -p 5433 -U postgres -d ecommerce_db -f migrations/001_init.sql
```

### 12.6 Code Organization Guidelines

#### Adding New Entity

1. **Create Model** (`internal/models/entity.go`)

   ```go
   type Entity struct {
       ID        uuid.UUID `json:"id"`
       // ... fields
   }
   ```

2. **Create Repository** (`internal/repository/entity_repo.go`)

   ```go
   type EntityRepository interface {
       Create(ctx context.Context, entity *models.Entity) error
       GetByID(ctx context.Context, id uuid.UUID) (*models.Entity, error)
   }
   ```

3. **Create Service** (`internal/service/entity_service.go`)

   ```go
   type EntityService interface {
       CreateEntity(ctx context.Context, req models.EntityRequest) (*models.Entity, error)
   }
   ```

4. **Create Handler** (`internal/handlers/entity.go`)

   ```go
   type EntityHandler struct {
       service service.EntityService
   }
   ```

5. **Register Routes** (`cmd/server/main.go`)

### 12.7 Common Development Tasks

#### Add New API Endpoint

1. Define request/response models in `models/`
2. Add validation tags
3. Implement business logic in service
4. Create handler function
5. Register route in `main.go`
6. Update OpenAPI specification

#### Add New Middleware

1. Create middleware file in `internal/middleware/`
2. Implement as `gin.HandlerFunc`
3. Register in middleware chain in `main.go`

#### Add Database Migration

1. Create new SQL file: `migrations/002_description.sql`
2. Add migration logic
3. Update Docker Compose volumes if needed

---

## 13. API Response Examples

### 13.1 Successful Responses

**User Registration:**

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "customer@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "customer",
      "created_at": "2026-02-07T10:00:00Z",
      "updated_at": "2026-02-07T10:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Get Products (Paginated):**

```json
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": {
    "products": [
      {
        "id": "7f8e9d10-a1b2-4c5d-9e8f-7a6b5c4d3e2f",
        "sku": "LAP-001",
        "name": "Gaming Laptop",
        "description": "High-performance gaming laptop",
        "price": 1499.99,
        "stock": 25,
        "category": "Electronics",
        "image_url": "https://...",
        "created_at": "2026-02-01T08:00:00Z",
        "updated_at": "2026-02-01T08:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 4,
      "total_pages": 1
    }
  }
}
```

**Cart with Items:**

```json
{
  "success": true,
  "message": "Cart retrieved successfully",
  "data": {
    "id": "a1b2c3d4-e5f6-7g8h-9i0j-k1l2m3n4o5p6",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "items": [
      {
        "id": "item-uuid",
        "cart_id": "cart-uuid",
        "product_id": "product-uuid",
        "product": {
          "id": "product-uuid",
          "name": "Gaming Laptop",
          "price": 1499.99,
          "stock": 25
        },
        "quantity": 2,
        "created_at": "2026-02-07T10:15:00Z"
      }
    ],
    "created_at": "2026-02-07T10:00:00Z",
    "updated_at": "2026-02-07T10:15:00Z"
  }
}
```

**Created Order:**

```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "id": "order-uuid",
    "user_id": "user-uuid",
    "order_number": "ORD-1707302400-a1b2c3d4",
    "total_amount": 2999.98,
    "status": "pending",
    "payment_method": "cc",
    "shipping_address": {
      "full_name": "John Doe",
      "street": "123 Main St",
      "city": "New York",
      "state": "NY",
      "country": "USA",
      "postal_code": "10001",
      "phone": "+1234567890"
    },
    "billing_address": {
      /* same structure */
    },
    "items": [
      {
        "id": "item-uuid",
        "order_id": "order-uuid",
        "product_id": "product-uuid",
        "product": {
          /* product details */
        },
        "quantity": 2,
        "price_at_time": 1499.99,
        "created_at": "2026-02-07T10:20:00Z"
      }
    ],
    "created_at": "2026-02-07T10:20:00Z",
    "updated_at": "2026-02-07T10:20:00Z"
  }
}
```

### 13.2 Error Responses

**Validation Error:**

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": [
    {
      "field": "email",
      "message": "Invalid email format"
    },
    {
      "field": "password",
      "message": "Password must be at least 6 characters"
    }
  ]
}
```

**Authentication Error:**

```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "Invalid token"
}
```

**Authorization Error:**

```json
{
  "success": false,
  "message": "Forbidden",
  "error": "Admin access required"
}
```

**Business Logic Error:**

```json
{
  "success": false,
  "message": "Operation failed",
  "error": "Insufficient stock"
}
```

**Not Found Error:**

```json
{
  "success": false,
  "message": "Resource not found",
  "error": "Product not found"
}
```

---

## 14. Performance Considerations

### 14.1 Database Optimization

**Connection Pooling:**

```go
poolConfig.MaxConns = 25
poolConfig.MinConns = 5
poolConfig.MaxConnLifetime = time.Hour
poolConfig.MaxConnIdleTime = 30 * time.Minute
```

**Indexes:**

- Primary keys (automatic)
- Foreign keys
- `products.sku` (unique lookups)
- `products.category` (filtering)
- `orders.user_id` (user order history)
- `orders.status` (admin filtering)
- `stock_reservations.expires_at` (cleanup queries)

**Query Optimization:**

- Parameterized queries (prevent SQL injection + query plan caching)
- JOIN queries for related data (avoid N+1 queries)
- Aggregate queries for analytics
- LIMIT/OFFSET for pagination

### 14.2 Stock Reservation Cleanup

**Automatic Expiry:**
Stock reservations automatically expire based on `expires_at` timestamp. Queries always filter with `AND expires_at > NOW()`.

**Background Cleanup (Recommended for Production):**

```sql
-- Run periodically via cron job or background worker
DELETE FROM stock_reservations
WHERE expires_at < NOW();
```

### 14.3 Caching Strategy (Future Enhancement)

**Recommended Caching:**

1. **Product Catalog:** Cache product listings (TTL: 5 minutes)
2. **User Sessions:** Cache JWT validation results
3. **Analytics:** Cache admin dashboard metrics (TTL: 1 hour)

**Not Recommended to Cache:**

- Cart contents (real-time updates needed)
- Stock levels (accuracy critical)
- Order details (consistency critical)

---

## 15. API Testing with Postman

The project includes a Postman collection: `postman_collection.json`

**Import into Postman:**

1. Open Postman
2. File → Import
3. Select `postman_collection.json`

**Collection Variables:**

- `baseUrl`: `http://localhost:8080`
- `token`: (Auto-set by login request)

**Test Sequence:**

1. Health Check
2. Register User
3. Login
4. Get Products
5. Add to Cart
6. Get Cart
7. Create Order
8. Get User Orders
9. Admin: Update Order Status
10. Create Return Request

---

## 16. Future Enhancements

### Recommended Features

1. **Email Notifications**
   - Order confirmation
   - Shipping updates
   - Return status updates

2. **Payment Gateway Integration**
   - Stripe/PayPal integration
   - Real payment processing
   - Webhook handling

3. **Real-time Updates**
   - WebSocket support
   - Order status notifications
   - Stock level alerts

4. **Advanced Search**
   - Full-text search
   - Filters (price range, ratings)
   - Sorting options

5. **Product Reviews & Ratings**
   - Customer reviews
   - Rating system
   - Review moderation

6. **Inventory Management**
   - Low stock alerts
   - Automatic reordering
   - Supplier management

7. **Discount & Coupons**
   - Promo code system
   - Percentage/fixed discounts
   - Expiry dates

8. **Shipping Integration**
   - Multiple shipping providers
   - Real-time tracking
   - Shipping cost calculation

9. **Analytics Dashboard**
   - Sales reports
   - Customer insights
   - Product performance

10. **API Rate Limiting**
    - Per-user rate limits
    - IP-based throttling
    - DDoS protection

---

## 17. Troubleshooting

### Common Issues

**Database Connection Failed:**

```
Solution:
1. Check PostgreSQL is running: docker-compose ps
2. Verify credentials in .env file
3. Check port 5433 is not in use
4. Wait for health check to pass
```

**JWT Token Invalid:**

```
Solution:
1. Check JWT_SECRET matches between services
2. Verify token hasn't expired (24h default)
3. Ensure Authorization header format: Bearer <token>
```

**Stock Reservation Issues:**

```
Solution:
1. Check stock_reservations table for expired entries
2. Run cleanup query manually
3. Adjust STOCK_RESERVATION_TTL_MINUTES if needed
```

**Order Creation Fails:**

```
Solution:
1. Validate cart is not empty
2. Check stock availability
3. Verify addresses are complete
4. Check database transaction logs
```

---

## 18. Contributing Guidelines

### Code Style

- Follow Go standard formatting (`gofmt`)
- Use meaningful variable names
- Add comments for complex logic
- Keep functions small and focused

### Commit Messages

```
feat: Add product search functionality
fix: Resolve cart quantity update issue
docs: Update API documentation
refactor: Simplify order creation logic
test: Add unit tests for auth service
```

###Testing

- Write unit tests for services
- Test edge cases
- Maintain test coverage > 70%
- Test error scenarios

---

## 19. License & Contact

**Project:** E-Commerce Backend API  
**Author:** Jitender  
**Version:** 1.0.0  
**Date:** February 2026

For questions, issues, or contributions, please refer to the project repository.

---

**End of Documentation**

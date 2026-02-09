# E-Commerce Backend - Quick Reference Guide

## Quick Start

### Start Application (Docker)

```bash
cd ecommerce-backend-go
docker-compose up -d
```

### Access

- **API Base URL:** `http://localhost:8080`
- **Health Check:** `http://localhost:8080/health`
- **API Version:** `http://localhost:8080/api/v1`

### Default Admin Account

```
Email: admin@example.com
Password: admin123
```

---

## API Endpoints Cheat Sheet

### Public Endpoints (No Auth Required)

| Method | Endpoint                | Description         |
| ------ | ----------------------- | ------------------- |
| GET    | `/health`               | Health check        |
| GET    | `/ready`                | Readiness check     |
| POST   | `/api/v1/auth/register` | Register new user   |
| POST   | `/api/v1/auth/login`    | Login user          |
| POST   | `/api/v1/auth/refresh`  | Refresh token       |
| GET    | `/api/v1/products`      | List products       |
| GET    | `/api/v1/products/:id`  | Get product details |

### Protected Endpoints (Require Auth)

#### User Profile

```
GET    /api/v1/users/profile          - Get profile
PUT    /api/v1/users/profile          - Update profile
PUT    /api/v1/users/change-password  - Change password
```

#### Shopping Cart

```
GET    /api/v1/cart                - Get cart
POST   /api/v1/cart/items          - Add to cart
PUT    /api/v1/cart/items/:itemId  - Update quantity
DELETE /api/v1/cart/items/:itemId  - Remove item
DELETE /api/v1/cart                - Clear cart
GET    /api/v1/cart/validate       - Validate cart
```

#### Orders

```
POST   /api/v1/orders              - Create order
GET    /api/v1/orders              - Get user orders
GET    /api/v1/orders/:id          - Get order details
PUT    /api/v1/orders/:id/cancel   - Cancel order
GET    /api/v1/orders/:id/payment  - Get payment info
```

#### Payments

```
POST   /api/v1/payments              - Create payment
POST   /api/v1/payments/:id/verify   - Verify payment
```

#### Returns

```
POST   /api/v1/returns     - Create return request
GET    /api/v1/returns     - Get user returns
GET    /api/v1/returns/:id - Get return details
```

### Admin Endpoints (Require Admin Role)

#### Products

```
POST   /api/v1/admin/products      - Create product
GET    /api/v1/admin/products      - List all products
PUT    /api/v1/admin/products/:id  - Update product
DELETE /api/v1/admin/products/:id  - Delete product
GET    /api/v1/admin/products/top  - Top products
```

#### Orders

```
GET  /api/v1/admin/orders            - All orders
GET  /api/v1/admin/orders/recent     - Recent orders
GET  /api/v1/admin/orders/:id        - Order details
PUT  /api/v1/admin/orders/:id/status - Update status
GET  /api/v1/admin/analytics         - Analytics
```

#### Users

```
GET  /api/v1/admin/users         - All users
PUT  /api/v1/admin/users/:id/role - Update role
```

#### Returns

```
GET  /api/v1/admin/returns                   - All returns
POST /api/v1/admin/returns/:returnId/process - Process return
```

---

## Common Request Examples

### Register User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Response:**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {...},
    "access_token": "eyJhbGc..."
  }
}
```

### List Products

```bash
curl http://localhost:8080/api/v1/products?page=1&limit=10
```

### Add to Cart

```bash
curl -X POST http://localhost:8080/api/v1/cart/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "product_id": "7f8e9d10-a1b2-4c5d-9e8f-7a6b5c4d3e2f",
    "quantity": 2
  }'
```

### Get Cart

```bash
curl http://localhost:8080/api/v1/cart \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Create Order

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
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
      "full_name": "John Doe",
      "street": "123 Main St",
      "city": "New York",
      "state": "NY",
      "country": "USA",
      "postal_code": "10001",
      "phone": "+1234567890"
    }
  }'
```

### Admin: Create Product

```bash
curl -X POST http://localhost:8080/api/v1/admin/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -d '{
    "sku": "NEW-001",
    "name": "New Product",
    "description": "Product description",
    "price": 99.99,
    "stock": 50,
    "category": "Electronics",
    "image_url": "https://example.com/image.jpg"
  }'
```

### Admin: Update Order Status

```bash
curl -X PUT http://localhost:8080/api/v1/admin/orders/ORDER_ID/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -d '{
    "status": "shipped"
  }'
```

---

## Query Parameters

### Pagination

```
?page=1&limit=10
```

### Product Filtering

```
?category=Electronics
?search=laptop
```

### Order Filtering (Admin)

```
?status=pending
?range_days=30
```

### User Filtering (Admin)

```
?range_days=7  (users registered in last 7 days)
```

---

## Response Format

### Success Response

```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    /* response data */
  }
}
```

### Error Response

```json
{
  "success": false,
  "message": "Operation failed",
  "error": "Error details"
}
```

### Validation Error

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": [{ "field": "email", "message": "Invalid email" }]
}
```

---

## Environment Variables

### Required

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ecommerce_db
JWT_SECRET=your-secret-key
```

### Optional

```bash
PORT=8080
ENV=development
JWT_EXPIRY_HOURS=24
ALLOWED_ORIGINS=http://localhost:3000
STOCK_RESERVATION_TTL_MINUTES=10
DB_SSLMODE=disable
```

---

## Docker Commands

### Start Services

```bash
docker-compose up -d
```

### View Logs

```bash
docker-compose logs -f
docker-compose logs -f backend
docker-compose logs -f postgres
```

### Stop Services

```bash
docker-compose down
```

### Restart Service

```bash
docker-compose restart backend
```

### Rebuild and Start

```bash
docker-compose up --build -d
```

### Remove Everything

```bash
docker-compose down -v --remove-orphans
```

### Check Status

```bash
docker-compose ps
```

---

## Database Commands

### Connect to Database

```bash
psql -h localhost -p 5433 -U postgres -d ecommerce_db
```

### Common SQL Queries

#### List All Tables

```sql
\dt
```

#### Count Users

```sql
SELECT COUNT(*) FROM users;
```

#### View Products

```sql
SELECT id, sku, name, price, stock_quantity FROM products;
```

#### View Recent Orders

```sql
SELECT id, order_number, total_amount, status, created_at
FROM orders
ORDER BY created_at DESC
LIMIT 10;
```

#### Check Stock Reservations

```sql
SELECT
    p.name,
    sr.quantity,
    sr.expires_at,
    CASE
        WHEN sr.expires_at > NOW() THEN 'Active'
        ELSE 'Expired'
    END as status
FROM stock_reservations sr
JOIN products p ON sr.product_id = p.id;
```

#### View Available Stock

```sql
SELECT
    p.name,
    p.stock_quantity as actual_stock,
    COALESCE(SUM(sr.quantity), 0) as reserved,
    p.stock_quantity - COALESCE(SUM(sr.quantity), 0) as available
FROM products p
LEFT JOIN stock_reservations sr
    ON p.id = sr.product_id
    AND sr.expires_at > NOW()
GROUP BY p.id, p.name, p.stock_quantity;
```

---

## Order Status Workflow

```
pending → processing → shipped → delivered → completed
   ↓           ↓
cancelled   cancelled
```

### Status Transitions (Admin)

- `pending` → `processing`, `cancelled`
- `processing` → `shipped`, `delivered`, `completed`, `cancelled`
- `shipped` → `delivered`, `completed`
- `delivered` → `completed`

### Payment Creation

- **CC/DC:** Immediately on order creation
- **COD:** When status changes to `delivered`

---

## Troubleshooting

### Database Connection Error

```bash
# Check PostgreSQL is running
docker-compose ps

# Restart PostgreSQL
docker-compose restart postgres

# Check logs
docker-compose logs postgres
```

### Port Already in Use

```bash
# Check what's using port 8080
netstat -ano | findstr :8080

# Kill process (Windows)
taskkill /PID <process_id> /F
```

### JWT Token Invalid

- Token expired (24h default)
- Wrong JWT_SECRET
- Malformed token

### Stock Reservation Issues

```sql
-- Clean up expired reservations
DELETE FROM stock_reservations
WHERE expires_at < NOW();

-- View all reservations
SELECT * FROM stock_reservations;
```

---

## Testing Workflow

### 1. Test Registration

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test123","first_name":"Test","last_name":"User"}'
```

### 2. Test Login (Save Token)

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test123"}'
```

### 3. Test Get Products

```bash
curl http://localhost:8080/api/v1/products
```

### 4. Test Add to Cart

```bash
curl -X POST http://localhost:8080/api/v1/cart/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"product_id":"PRODUCT_ID","quantity":1}'
```

### 5. Test Get Cart

```bash
curl http://localhost:8080/api/v1/cart \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 6. Test Create Order

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{...order payload...}'
```

---

## Performance Tips

### Database Indexes

- Already created on key columns
- Check query performance with `EXPLAIN ANALYZE`

### Connection Pool Settings

```go
MaxConns = 25
MinConns = 5
MaxConnLifetime = 1 hour
MaxConnIdleTime = 30 minutes
```

### Stock Reservation Cleanup

Run periodically:

```sql
DELETE FROM stock_reservations WHERE expires_at < NOW();
```

---

## Security Checklist

- ✅ JWT secret configured (not default)
- ✅ Database password changed
- ✅ CORS origins configured
- ✅ Rate limiting enabled (optional)
- ✅ HTTPS in production
- ✅ Environment variables secured
- ✅ Database backups configured

---

## File Structure Reference

```
ecommerce-backend-go/
├── cmd/server/main.go           # Entry point
├── internal/
│   ├── config/                  # Configuration
│   ├── handlers/                # HTTP handlers
│   ├── middleware/              # Middleware
│   ├── models/                  # Domain models
│   ├── repository/              # Data access
│   └── service/                 # Business logic
├── pkg/
│   ├── database/                # DB connection
│   └── utils/                   # Utilities
├── migrations/                  # SQL migrations
├── docker-compose.yml           # Docker config
├── Dockerfile                   # Container image
├── go.mod                       # Dependencies
├── openapi.yaml                 # API spec
└── postman_collection.json      # API tests
```

---

## Key Concepts

### Repository Pattern

Data access abstraction layer

### Service Layer

Business logic and orchestration

### Middleware

Cross-cutting concerns (auth, logging, etc.)

### JWT Authentication

Token-based stateless auth

### Stock Reservation

Prevents overselling with temporary locks

### Transaction Management

ACID compliance for complex operations

---

**Quick Tip:** Use the Postman collection (`postman_collection.json`) for easy API testing!

---

**End of Quick Reference**

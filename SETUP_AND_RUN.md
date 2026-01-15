# Ecommerce Backend GO - Setup & Run Guide

## Project Overview

- **Language**: Go 1.24
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL 15
- **Container**: Docker & Docker Compose
- **Port**: 8080 (Backend), 5433 (Database)

---

## Prerequisites

Before running the project, ensure you have:

- **Docker** installed and running
- **Docker Compose** installed
- **Go 1.24+** (if running locally without Docker)
- **PostgreSQL 15** (if running locally without Docker)
- **Git** (for cloning the repository)

### Install Docker & Docker Compose

```bash
# Windows (using Chocolatey)
choco install docker-desktop
choco install docker-compose

# Or download from:
# https://www.docker.com/products/docker-desktop
# https://docs.docker.com/compose/install/
```

### Install Go (Optional - for local development)

```bash
# Download from https://golang.org/dl/
# Or using package manager:
# Windows (Chocolatey)
choco install golang

# Verify installation
go version
```

---

## Quick Start (Docker Compose - Recommended)

### 1. **Clone/Open Project**

```bash
cd d:\jitender-personal\GO-Playlist\ecommerce-backend-go
```

### 2. **Start All Services**

```bash
# Start backend and database in foreground (see logs)
docker-compose up

# OR start in background
docker-compose up -d

# With rebuild (if code changed)
docker-compose up --build -d
```

### 3. **Verify Services Are Running**

```bash
# Check container status
docker-compose ps

# Expected output:
# NAME                 STATUS
# ecommerce-backend    Up 2 seconds
# ecommerce-db         Up 3 seconds

# Check if backend is responding
curl http://localhost:8080/health

# Expected response:
# {"status":"healthy","database":"connected",...}
```

### 4. **Stop All Services**

```bash
# Stop containers (keep data)
docker-compose down

# Stop and remove volumes (delete database data)
docker-compose down -v

# Stop and remove all (clean slate)
docker-compose down -v --remove-orphans
```

---

## Running Locally (Without Docker)

### 1. **Install Dependencies**

```bash
cd ecommerce-backend-go
go mod download
go mod tidy
```

### 2. **Setup PostgreSQL**

```bash
# Windows (using Chocolatey)
choco install postgresql

# Create database
createdb -U postgres ecommerce_db

# Connect and run migrations
psql -U postgres -d ecommerce_db -f migrations/001_init.sql
```

### 3. **Configure Environment Variables**

Create `.env` file in project root:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ecommerce_db
DB_SSLMODE=disable

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
ENVIRONMENT=development

# JWT Configuration
JWT_SECRET=your_super_secret_key_change_this
JWT_EXPIRATION=3600
```

### 4. **Build the Project**

```bash
cd cmd/server
go build -o server

# Or from root
go build -o ./bin/server ./cmd/server
```

### 5. **Run the Application**

```bash
# From project root
./cmd/server/server

# Or if built to bin
./bin/server

# Expected output:
# ✅ Successfully connected to PostgreSQL database
# [GIN-debug] Listening and serving HTTP on [::]:8080
```

### 6. **Test the API**

```bash
# Health check
curl http://localhost:8080/health

# Register user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Pass123","first_name":"Test","last_name":"User"}'
```

---

## Building & Deployment

### Build Docker Image Manually

```bash
# Build image
docker build -t ecommerce-backend:latest .

# Run container
docker run -d \
  --name ecommerce-backend \
  -p 8080:8080 \
  --env-file .env \
  ecommerce-backend:latest
```

### Build Go Binary for Production

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o server cmd/server/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o server.exe cmd/server/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o server cmd/server/main.go
```

---

## Database Operations

### View Database Migrations

```bash
# From project root
cat migrations/001_init.sql

# Or list all migrations
ls migrations/
```

### Run Migrations Manually

```bash
# Using psql (local PostgreSQL)
psql -U postgres -d ecommerce_db -f migrations/001_init.sql

# Using Docker exec (Docker database)
docker exec ecommerce-db psql -U postgres -d ecommerce_db -f /migrations/001_init.sql
```

### Connect to Database Directly

```bash
# Local PostgreSQL
psql -U postgres -d ecommerce_db

# Docker PostgreSQL
docker exec -it ecommerce-db psql -U postgres -d ecommerce_db

# Example queries:
# \dt                          -- List all tables
# SELECT * FROM users;        -- View users
# SELECT * FROM products;     -- View products
# SELECT * FROM orders;       -- View orders
# \q                          -- Quit psql
```

### Backup Database

```bash
# Create backup
docker exec ecommerce-db pg_dump -U postgres ecommerce_db > backup.sql

# Restore backup
docker exec -i ecommerce-db psql -U postgres ecommerce_db < backup.sql
```

---

## Development Commands

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/service/...

# Run with coverage
go test -cover ./...
```

### Build & Run with Make (if Makefile exists)

```bash
# View available targets
make help

# Build
make build

# Run
make run

# Test
make test

# Clean
make clean
```

### Format Code

```bash
# Format all Go files
go fmt ./...

# Use gofmt on specific directory
gofmt -w internal/

# Run golint (if installed)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run ./...
```

### View Logs

```bash
# Docker container logs
docker logs ecommerce-backend

# Follow logs (real-time)
docker logs -f ecommerce-backend

# Last 50 lines
docker logs --tail 50 ecommerce-backend

# Database logs
docker logs ecommerce-db
```

---

## Testing APIs with cURL

### User Registration

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email":"user@example.com",
    "password":"Password123",
    "first_name":"John",
    "last_name":"Doe"
  }'
```

### User Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email":"user@example.com",
    "password":"Password123"
  }'
```

### Get Products

```bash
curl -X GET http://localhost:8080/api/v1/products
```

### Add to Cart

```bash
curl -X POST http://localhost:8080/api/v1/cart/items \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id":"06c2a670-8ae9-4f39-aa0f-ad818073bb3b",
    "quantity":2
  }'
```

### Get Cart

```bash
curl -X GET http://localhost:8080/api/v1/cart \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Create Order

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "shipping_address":{
      "full_name":"John Doe",
      "street":"123 Main St",
      "city":"New York",
      "state":"NY",
      "country":"USA",
      "postal_code":"10001"
    },
    "billing_address":{
      "full_name":"John Doe",
      "street":"123 Main St",
      "city":"New York",
      "state":"NY",
      "country":"USA",
      "postal_code":"10001"
    }
  }'
```

### Get Orders

```bash
curl -X GET http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## Testing APIs with PowerShell

### Health Check

```powershell
$health = Invoke-WebRequest -Uri "http://localhost:8080/health" -Method GET -UseBasicParsing
$health.Content | ConvertFrom-Json | Format-List
```

### Register User

```powershell
$body = @{
    email = "test@example.com"
    password = "Pass123"
    first_name = "Test"
    last_name = "User"
} | ConvertTo-Json

$response = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/register" `
  -Method POST `
  -ContentType "application/json" `
  -Body $body `
  -UseBasicParsing

$response.Content | ConvertFrom-Json | ConvertTo-Json -Depth 3
```

### Get Products

```powershell
$response = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/products" `
  -Method GET `
  -UseBasicParsing

($response.Content | ConvertFrom-Json).data.products | Format-Table Name, Price, Stock
```

---

## Troubleshooting

### Port Already in Use

```bash
# Find process using port 8080
netstat -ano | findstr :8080

# Kill process (Windows)
taskkill /PID <PID> /F

# Or change port in docker-compose.yml or .env
```

### Database Connection Failed

```bash
# Check if PostgreSQL is running
docker ps | grep ecommerce-db

# Check database logs
docker logs ecommerce-db

# Verify connection string in .env
# Format: postgres://user:password@host:port/database
```

### Container Won't Start

```bash
# Check container logs
docker logs ecommerce-backend

# Remove and rebuild
docker-compose down -v
docker-compose up --build

# Check Docker daemon is running
docker ps
```

### Go Modules Issues

```bash
# Clear cache
go clean -modcache

# Re-download dependencies
go mod download

# Verify dependencies
go mod verify

# Tidy modules
go mod tidy
```

---

## Project Structure

```
ecommerce-backend-go/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── config/                  # Configuration
│   ├── handlers/                # API handlers
│   ├── middleware/              # Middleware (auth, CORS, logging)
│   ├── models/                  # Data models
│   ├── repository/              # Database operations
│   ├── service/                 # Business logic
├── migrations/
│   └── 001_init.sql            # Database schema
├── pkg/
│   ├── database/               # Database connection
│   └── utils/                  # Utilities (JWT, response formatting)
├── docker-compose.yml          # Docker Compose configuration
├── Dockerfile                  # Docker image definition
├── go.mod                      # Go dependencies
├── go.sum                      # Go dependency checksums
├── .env                        # Environment variables
└── README.md                   # Project documentation
```

---

## API Endpoints Summary

### Authentication

- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user

### Users

- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update profile
- `PUT /api/v1/users/change-password` - Change password

### Products

- `GET /api/v1/products` - List products (paginated)
- `GET /api/v1/products/:id` - Get single product

### Cart

- `GET /api/v1/cart` - Get cart
- `POST /api/v1/cart/items` - Add to cart
- `PUT /api/v1/cart/items/:itemId` - Update cart item
- `DELETE /api/v1/cart/items/:itemId` - Remove from cart
- `DELETE /api/v1/cart` - Clear cart
- `GET /api/v1/cart/validate` - Validate cart

### Orders

- `POST /api/v1/orders` - Create order
- `GET /api/v1/orders` - Get user orders
- `GET /api/v1/orders/:id` - Get order details
- `PUT /api/v1/orders/:id/cancel` - Cancel order

### Health

- `GET /health` - Health check

---

## Environment Variables

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ecommerce_db
DB_SSLMODE=disable

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
ENVIRONMENT=development

# JWT
JWT_SECRET=your_secret_key
JWT_EXPIRATION=3600
```

---

## Useful Links

- [Go Documentation](https://golang.org/doc/)
- [Gin Web Framework](https://gin-gonic.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)

---

## Support

For issues or questions:

1. Check logs: `docker logs ecommerce-backend`
2. Verify database: `docker exec -it ecommerce-db psql -U postgres -d ecommerce_db`
3. Check API: `curl http://localhost:8080/health`
4. Review .env configuration

---

**Last Updated**: January 15, 2026
**Project Version**: 1.0.0

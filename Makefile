.PHONY: all build run test clean migrate db-up db-down

# Variables
APP_NAME = ecommerce-backend
BINARY_NAME = bin/$(APP_NAME)
MIGRATIONS_DIR = migrations

all: build

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BINARY_NAME) cmd/server/main.go
	@echo "Build complete: $(BINARY_NAME)"

run: build
	@echo "Starting $(APP_NAME)..."
	@./$(BINARY_NAME)

dev:
	@echo "Starting in development mode..."
	@go run cmd/server/main.go

test:
	@echo "Running tests..."
	@go test ./... -v

clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@rm -rf coverage.out
	@go clean

# Database commands
db-up:
	@echo "Starting PostgreSQL database..."
	@docker-compose up -d postgres

db-down:
	@echo "Stopping PostgreSQL database..."
	@docker-compose down

migrate:
	@echo "Running migrations..."
	@psql -h localhost -U postgres -d ecommerce_db -f $(MIGRATIONS_DIR)/001_init.sql

migrate-create:
	@read -p "Enter migration name: " name; \
	mkdir -p $(MIGRATIONS_DIR); \
	timestamp=$$(date +%Y%m%d%H%M%S); \
	touch $(MIGRATIONS_DIR)/$${timestamp}_$${name}.sql

# Code quality
lint:
	@echo "Running linter..."
	@golangci-lint run

fmt:
	@echo "Formatting code..."
	@gofmt -w .

vet:
	@echo "Running vet..."
	@go vet ./...

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):latest .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

docker-compose-up:
	@echo "Starting with Docker Compose..."
	@docker-compose up --build

# Help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  run            - Build and run the application"
	@echo "  dev            - Run in development mode"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  db-up          - Start PostgreSQL database"
	@echo "  db-down        - Stop PostgreSQL database"
	@echo "  migrate        - Run database migrations"
	@echo "  migrate-create - Create a new migration file"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-compose-up - Start with Docker Compose"
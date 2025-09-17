.PHONY: help build run test test-coverage clean docker-build docker-run docker-compose-up docker-compose-down migrate-up migrate-down lint

# Default target
help:
	@echo "Available targets:"
	@echo "  build              - Build the application"
	@echo "  run                - Run the application locally"
	@echo "  test               - Run unit tests"
	@echo "  test-coverage      - Run tests with coverage report"
	@echo "  clean              - Clean build artifacts"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run Docker container"
	@echo "  docker-compose-up  - Start all services with Docker Compose"
	@echo "  docker-compose-down- Stop all services"
	@echo "  migrate-up         - Run database migrations"
	@echo "  migrate-down       - Rollback database migrations"
	@echo "  lint               - Run linter"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/server cmd/oms/main.go

# Run the application locally
run:
	@echo "Running application..."
	docker-compose up -d --build

# go run cmd/oms/main.go


# Run unit tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t squadio-oms-ai:latest .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --name oms-app squadio-oms-ai:latest

# Start all services with Docker Compose
docker-compose-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

# Stop all services
docker-compose-down:
	@echo "Stopping services..."
	docker-compose down

# Run database migrations
migrate-up:
	@echo "Running database migrations..."
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/squadio?sslmode=disable" up

# Rollback database migrations
migrate-down:
	@echo "Rolling back database migrations..."
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/squadio?sslmode=disable" down

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Setup development environment
setup: install-tools deps
	@echo "Setting up development environment..."
	@echo "Please start PostgreSQL and run 'make migrate-up'"

# Full test suite (unit + integration)
test-all: test
	@echo "Running full test suite..."

# Development workflow
dev: clean build test
	@echo "Development workflow completed"

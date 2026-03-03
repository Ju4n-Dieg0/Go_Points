.PHONY: help build run test clean docker-build docker-up docker-down migrate dev install lint

# Variables
APP_NAME=go-points-api
BINARY_NAME=server
MAIN_PATH=./cmd/server
DOCKER_COMPOSE=docker-compose

# Default target
help:
	@echo "Available commands:"
	@echo "  make install       - Install dependencies"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make dev           - Run the application with auto-reload"
	@echo "  make test          - Run tests"
	@echo "  make lint          - Run linters"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start services with Docker Compose"
	@echo "  make docker-down   - Stop services with Docker Compose"
	@echo "  make migrate       - Run database migrations"

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Build the application
build:
	@echo "Building application..."
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run: build
	@echo "Running application..."
	./bin/$(BINARY_NAME)

# Run with air for development (auto-reload)
dev:
	@echo "Running in development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Install with: go install github.com/air-verse/air@latest"; \
		echo "Running without hot reload..."; \
		go run $(MAIN_PATH)/main.go; \
	fi

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linters
lint:
	@echo "Running linters..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "golangci-lint not installed. Install from: https://golangci-lint.run/usage/install/"; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):latest .

# Start Docker Compose services
docker-up:
	@echo "Starting services..."
	$(DOCKER_COMPOSE) up -d

# Stop Docker Compose services
docker-down:
	@echo "Stopping services..."
	$(DOCKER_COMPOSE) down

# Restart Docker Compose services
docker-restart:
	@echo "Restarting services..."
	$(DOCKER_COMPOSE) restart

# View Docker Compose logs
docker-logs:
	$(DOCKER_COMPOSE) logs -f

# Run database migrations (when implemented)
migrate:
	@echo "Running migrations..."
	go run $(MAIN_PATH)/main.go

# Generate mocks (when implemented)
mock:
	@echo "Generating mocks..."
	@if command -v mockgen > /dev/null; then \
		go generate ./...; \
	else \
		echo "mockgen not installed. Install with: go install github.com/golang/mock/mockgen@latest"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...

# Run all checks
check: fmt vet lint test
	@echo "All checks passed!"

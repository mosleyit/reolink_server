.PHONY: help build run test clean deps docker-build docker-up docker-down migrate-up migrate-down

# Variables
APP_NAME=reolink_server
BINARY_NAME=reolink-server
VERSION?=1.0.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) cmd/server/main.go

run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	go run $(LDFLAGS) cmd/server/main.go

dev: ## Run with hot reload (requires air)
	@echo "Running in development mode..."
	air

test: ## Run tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):$(VERSION) .
	docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest

docker-up: ## Start Docker Compose services
	@echo "Starting Docker Compose services..."
	docker-compose up -d

docker-down: ## Stop Docker Compose services
	@echo "Stopping Docker Compose services..."
	docker-compose down

docker-logs: ## View Docker Compose logs
	docker-compose logs -f

migrate-create: ## Create a new migration (usage: make migrate-create NAME=migration_name)
	@echo "Creating migration: $(NAME)"
	@mkdir -p migrations
	@touch migrations/$(shell date +%Y%m%d%H%M%S)_$(NAME).up.sql
	@touch migrations/$(shell date +%Y%m%d%H%M%S)_$(NAME).down.sql

migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	# TODO: Add migration tool command

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	# TODO: Add migration tool command

.DEFAULT_GOAL := help


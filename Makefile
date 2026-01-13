.PHONY: help build run test clean docker-up docker-down lint fmt

help: ## Display this help screen
	@echo "Available commands:"
	@grep -h "^[a-zA-Z_-]*:.*##" $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "} {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the server binary
	@echo "Building server..."
	@go build -v -o bin/server ./cmd/server/
	@echo "✓ Build complete: bin/server"

run: build ## Build and run the server
	@echo "Starting server..."
	@./bin/server

test: ## Run all tests
	@echo "Running tests..."
	@go test -v -race -cover ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "✓ Coverage report: coverage.html"

bench: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@gofmt -w -s .
	@echo "✓ Code formatted"

fmt-check: ## Check code formatting
	@echo "Checking code formatting..."
	@gofmt -l .

clean: ## Remove build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean
	@echo "✓ Clean complete"

docker-up: ## Start PostgreSQL container
	@echo "Starting PostgreSQL..."
	@docker compose up -d
	@echo "✓ PostgreSQL is running"

docker-down: ## Stop PostgreSQL container
	@echo "Stopping PostgreSQL..."
	@docker compose down
	@echo "✓ PostgreSQL stopped"

docker-logs: ## View PostgreSQL logs
	@docker compose logs -f postgres

migrate: ## Initialize database schema
	@echo "Initializing database..."
	@go run ./cmd/server/main.go --migrate-only

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@echo "✓ Dependencies downloaded"

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "✓ Dependencies tidied"

dev: docker-up ## Start development environment (database + server)
	@echo "Starting development environment..."
	@sleep 2
	@go run ./cmd/server/main.go

all: deps lint test build ## Run all checks and build
	@echo "✓ All checks passed and build complete"

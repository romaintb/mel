.PHONY: help build test clean fmt lint vet coverage bench install-tools

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build the application
build: ## Build the application
	@echo "Building mel..."
	@go build -o bin/mel ./cmd/mel

# Build with version information
build-release: ## Build with version information
	@echo "Building mel with version info..."
	@go build -ldflags="-X main.version=$(shell git describe --tags --always --dirty)" -o bin/mel ./cmd/mel

# Install development tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install mvdan.cc/gofumpt@latest
	@go install golang.org/x/tools/cmd/goimports@latest

# Run tests
test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
bench: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Format code
fmt: ## Format Go code
	@echo "Formatting code..."
	@gofumpt -w .
	@goimports -w .

# Run linter
lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

# Run go vet
vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Run all checks
check: fmt lint vet test ## Run all checks (fmt, lint, vet, test)

# Development mode with file watching
dev: ## Run in development mode (requires air)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Or run: make build && ./bin/mel"; \
	fi

# Install the application
install: build ## Install the application
	@echo "Installing mel..."
	@cp bin/mel /usr/local/bin/

# Uninstall the application
uninstall: ## Uninstall the application
	@echo "Uninstalling mel..."
	@rm -f /usr/local/bin/mel

# Generate documentation
docs: ## Generate documentation
	@echo "Generating documentation..."
	@go doc -all ./...

# Show dependencies
deps: ## Show dependencies
	@echo "Dependencies:"
	@go list -m all

# Update dependencies
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

# Show module info
module-info: ## Show module info
	@echo "Module info:"
	@go mod graph

# Run with race detection
test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	@go test -race ./...

# Run specific test
test-file: ## Run tests in a specific file (usage: make test-file FILE=path/to/file_test.go)
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make test-file FILE=path/to/file_test.go"; \
		exit 1; \
	fi
	@echo "Running tests in $(FILE)..."
	@go test -v $(FILE)

# Show help for a specific target
help-target: ## Show help for a specific target (usage: make help-target TARGET=target_name)
	@if [ -z "$(TARGET)" ]; then \
		echo "Usage: make help-target TARGET=target_name"; \
		exit 1; \
	fi
	@grep -E '^$(TARGET):.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' 
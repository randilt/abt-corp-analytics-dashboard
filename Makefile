.PHONY: build test test-unit test-integration test-coverage clean run dev docker help preprocess

# Variables
BINARY_NAME=analytics-dashboard
MAIN_PATH=./cmd/server
BUILD_DIR=./bin
TEST_COVERAGE_DIR=./coverage
GO_FILES=$(shell find . -name "*.go" -type f -not -path "./vendor/*")

# Default target
all: clean test build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Development mode with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@echo "Starting development server with hot reload..."
	air -c .air.toml

# Run all tests
test: test-unit test-integration

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	go test -v -race -timeout 30s ./tests/unit/...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v -race -timeout 60s ./tests/integration/...

# Generate test coverage report
test-coverage:
	@echo "Generating test coverage report..."
	@mkdir -p $(TEST_COVERAGE_DIR)
	go test -v -race -coverprofile=$(TEST_COVERAGE_DIR)/coverage.out ./...
	go tool cover -html=$(TEST_COVERAGE_DIR)/coverage.out -o $(TEST_COVERAGE_DIR)/coverage.html
	go tool cover -func=$(TEST_COVERAGE_DIR)/coverage.out
	@echo "Coverage report generated: $(TEST_COVERAGE_DIR)/coverage.html"

# Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Preprocess data (run the preprocessing script)
preprocess:
	@echo "Preprocessing CSV data..."
	go run ./scripts/preprocess.go -csv ./data/raw/transactions.csv -cache ./data/processed/analytics_cache.json

# Load test (requires hey: go install github.com/rakyll/hey@latest)
load-test:
	@echo "Running load test..."
	hey -n 1000 -c 10 -m GET http://localhost:8080/api/v1/analytics

# Format Go code
fmt:
	@echo "Formatting Go code..."
	gofmt -s -w $(GO_FILES)

# Lint Go code (requires golangci-lint)
lint:
	@echo "Linting Go code..."
	golangci-lint run

# Vet Go code
vet:
	@echo "Vetting Go code..."
	go vet ./...

# Security scan (requires gosec: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
security:
	@echo "Running security scan..."
	gosec ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -rf $(TEST_COVERAGE_DIR)
	go clean -cache
	go clean -testcache

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Build Docker image
docker:
	@echo "Building Docker image..."
	docker build -t analytics-dashboard:latest -f deployments/docker/Dockerfile .

# Run with Docker
docker-run:
	@echo "Running with Docker..."
	docker run -p 8080:8080 -v $(PWD)/data:/app/data analytics-dashboard:latest

# Database migrations (for future use)
migrate-up:
	@echo "Running database migrations..."
	# migrate -path ./data/migrations -database "postgres://user:pass@localhost/db?sslmode=disable" up

migrate-down:
	@echo "Reverting database migrations..."
	# migrate -path ./data/migrations -database "postgres://user:pass@localhost/db?sslmode=disable" down

# Generate API documentation
docs:
	@echo "Generating API documentation..."
	# swag init -g cmd/server/main.go -o ./docs/swagger

# Performance profiling
profile:
	@echo "Starting CPU profiling..."
	go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./tests/benchmark/
	go tool pprof cpu.prof

# Memory profiling
profile-mem:
	@echo "Starting memory profiling..."
	go test -memprofile=mem.prof -bench=. ./tests/benchmark/
	go tool pprof mem.prof

# Check for vulnerabilities
vuln-check:
	@echo "Checking for vulnerabilities..."
	go list -json -m all | nancy sleuth

# Generate mocks (requires mockgen: go install github.com/golang/mock/mockgen@latest)
mocks:
	@echo "Generating mocks..."
	mockgen -source=internal/services/interfaces.go -destination=tests/mocks/services.go
	mockgen -source=internal/repository/interfaces.go -destination=tests/mocks/repository.go

# Setup development environment
setup-dev: deps
	@echo "Setting up development environment..."
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/rakyll/hey@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install github.com/golang/mock/mockgen@latest
	@echo "Creating data directories..."
	mkdir -p ./data/raw ./data/processed ./data/migrations
	@echo "Development environment setup complete!"

# Help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  run            - Build and run the application"
	@echo "  dev            - Run in development mode with hot reload"
	@echo "  test           - Run all tests"
	@echo "  test-unit      - Run unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-coverage  - Generate test coverage report"
	@echo "  benchmark      - Run benchmarks"
	@echo "  preprocess     - Preprocess CSV data"
	@echo "  load-test      - Run load test"
	@echo "  fmt            - Format Go code"
	@echo "  lint           - Lint Go code"
	@echo "  vet            - Vet Go code"
	@echo "  security       - Run security scan"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Install dependencies"
	@echo "  docker         - Build Docker image"
	@echo "  docker-run     - Run with Docker"
	@echo "  setup-dev      - Setup development environment"
	@echo "  help           - Show this help message"
.PHONY: all test unit-test integration-test lint build run-redis clean

# Project information
PROJECT_NAME := go-toolkit
MAIN_PACKAGE := ./cmd/examples/ratelimit_demo
BINARY_NAME := ratelimit-demo

# Tools
GO := go
DOCKER := docker
DOCKER_COMPOSE := docker-compose
GOLANGCI_LINT := golangci-lint

# Testing
TEST_FLAGS := -v
TEST_COVERAGE_FLAGS := -cover -coverprofile=coverage.out
UNIT_TEST_PATTERN := ".*_test\.go"
INTEGRATION_TEST_PATTERN := ".*_integration_test\.go"

# Docker
DOCKER_COMPOSE_FILE := ./docker/docker-compose.yml

# Default target
all: lint test build

# Build application
build:
	@echo "Building application..."
	$(GO) build -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)

# Run application
run: build
	@echo "Running application..."
	./bin/$(BINARY_NAME)

# Run all tests
test: unit-test integration-test

# Run unit tests
unit-test:
	@echo "Running unit tests..."
	$(GO) test $(TEST_FLAGS) $(TEST_COVERAGE_FLAGS) ./pkg/ratelimit/...

# Run integration tests
integration-test: run-redis
	@echo "Running integration tests..."
	$(GO) test $(TEST_FLAGS) -tags=integration ./test/ratelimit/...

# Run code quality checks
lint:
	@echo "Running linter..."
	$(GOLANGCI_LINT) run ./...

# Start Docker containers
docker-up:
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d

# Stop Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down

# Start only Redis container
run-redis:
	@echo "Starting Redis container..."
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d redis

# Generate test coverage report
coverage: test
	@echo "Generating test coverage report..."
	$(GO) tool cover -html=coverage.out

# Clean temporary files and build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin
	rm -f coverage.out

# Show help information
help:
	@echo "Available commands:"
	@echo "  make all            - Run lint, test and build"
	@echo "  make build          - Build application"
	@echo "  make run            - Build and run application"
	@echo "  make test           - Run all tests"
	@echo "  make unit-test      - Run only unit tests"
	@echo "  make integration-test - Run integration tests"
	@echo "  make lint           - Run code quality checks"
	@echo "  make docker-up      - Start all Docker containers"
	@echo "  make docker-down    - Stop all Docker containers"
	@echo "  make run-redis      - Start only Redis container"
	@echo "  make coverage       - Generate test coverage report"
	@echo "  make clean          - Clean temporary files and build artifacts"
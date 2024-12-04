PROJECT_NAME := docmind
MAIN_FILE := cmd/main.go
BINARY_NAME := docmind

GO := go
GOTEST := $(GO) test
GOVET := $(GO) vet
GOFMT := $(GO) fmt
GOMOD := $(GO) mod
GOGET := $(GO) get

COVERAGE_FILE := coverage.out
TEST_FLAGS := -v -race
TEST_TIMEOUT := 10m

DOCKER := docker
DOCKER_IMAGE := $(PROJECT_NAME)
DOCKER_TAG := latest

SHELL_SCRIPT := ./test-api.sh

.PHONY: all build clean test coverage lint fmt vet tidy help run test-api

all: lint test build

build:
	@echo "Building $(PROJECT_NAME)..."
	@$(GO) build -o $(BINARY_NAME) $(MAIN_FILE)

run:
	@echo "Running $(PROJECT_NAME)..."
	@$(GO) run $(MAIN_FILE)

clean:
	@echo "Cleaning..."
	@$(GO) clean
	@rm -f $(BINARY_NAME)
	@rm -f $(COVERAGE_FILE)
	@rm -f .last_token.json

test:
	@echo "Running tests..."
	@$(GOTEST) $(TEST_FLAGS) -timeout $(TEST_TIMEOUT) ./...

coverage:
	@echo "Running tests with coverage..."
	@$(GOTEST) -coverprofile=$(COVERAGE_FILE) ./...
	@$(GO) tool cover -html=$(COVERAGE_FILE)

benchmark:
	@echo "Running benchmark tests..."
	@$(GOTEST) -bench=. -benchmem ./...

fmt:
	@echo "Formatting code..."
	@$(GOFMT) ./...

tidy:
	@echo "Tidying dependencies..."
	@$(GOMOD) tidy



help:
	@echo "Available commands:"
	@echo "  make build              - Build the application"
	@echo "  make run               - Run the application"
	@echo "  make clean             - Clean build files"
	@echo "  make test              - Run all tests"
	@echo "  make coverage          - Generate test coverage report"
	@echo "  make benchmark         - Run benchmark tests"
	@echo "  make fmt               - Format code"
	@echo "  make tidy              - Tidy dependencies"
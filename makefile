# 项目变量
PROJECT_NAME := docmind
MAIN_FILE := cmd/main.go
BINARY_NAME := docmind

# Go 命令
GO := go
GOTEST := $(GO) test
GOVET := $(GO) vet
GOFMT := $(GO) fmt
GOMOD := $(GO) mod
GOGET := $(GO) get

# 测试相关
COVERAGE_FILE := coverage.out
TEST_FLAGS := -v -race
TEST_TIMEOUT := 10m

# Docker 相关
DOCKER := docker
DOCKER_COMPOSE := $(shell if command -v docker-compose >/dev/null 2>&1; then echo "docker-compose"; else echo "docker compose"; fi)
DOCKER_IMAGE := $(PROJECT_NAME)
DOCKER_TAG := latest
DOCKER_POSTGRES_SERVICE := postgres

# PostgreSQL 配置
POSTGRES_USER ?= docmind
POSTGRES_PASSWORD ?= docmind
POSTGRES_DB ?= docmind
POSTGRES_PORT ?= 5432

# 颜色定义
CYAN := \033[0;36m
GREEN := \033[0;32m
RED := \033[0;31m
NC := \033[0m # No Color

.PHONY: all build-go clean test coverage fmt tidy help run

# Go 相关命令
all: lint test build-go

build-go:
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

fmt:
	@echo "Formatting code..."
	@$(GOFMT) ./...

tidy:
	@echo "Tidying dependencies..."
	@$(GOMOD) tidy

# 检查 Docker 环境
.PHONY: check-docker
check-docker:
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "$(RED)Error: Docker is not installed$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Docker is installed$(NC)"
	@echo "$(CYAN)Using docker compose command: $(DOCKER_COMPOSE)$(NC)"

# Docker 和数据库相关命令
.PHONY: db-start
db-start: check-docker
	@echo "$(CYAN)Starting PostgreSQL container...$(NC)"
	@$(DOCKER_COMPOSE) up -d $(DOCKER_POSTGRES_SERVICE)
	@echo "$(GREEN)PostgreSQL is running on port $(POSTGRES_PORT)$(NC)"

.PHONY: db-stop
db-stop: check-docker
	@echo "$(CYAN)Stopping PostgreSQL container...$(NC)"
	@$(DOCKER_COMPOSE) stop $(DOCKER_POSTGRES_SERVICE)
	@echo "$(GREEN)PostgreSQL container stopped$(NC)"

.PHONY: db-restart
db-restart: db-stop db-start

.PHONY: db-logs
db-logs: check-docker
	@$(DOCKER_COMPOSE) logs -f $(DOCKER_POSTGRES_SERVICE)

.PHONY: db-shell
db-shell: check-docker
	@echo "$(CYAN)Connecting to PostgreSQL shell...$(NC)"
	@docker exec -it $$(docker ps -q -f name=postgres) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

.PHONY: db-status
db-status: check-docker
	@echo "$(CYAN)PostgreSQL container status:$(NC)"
	@docker ps -f name=postgres

.PHONY: db-clean
db-clean: check-docker
	@echo "$(RED)Warning: This will remove the PostgreSQL container and all its data$(NC)"
	@read -p "Are you sure? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo "$(CYAN)Removing PostgreSQL container and volumes...$(NC)"; \
		$(DOCKER_COMPOSE) down -v --remove-orphans; \
		echo "$(GREEN)Cleanup complete$(NC)"; \
	else \
		echo "$(CYAN)Cleanup cancelled$(NC)"; \
	fi

.PHONY: docker-build
docker-build: check-docker
	@echo "$(CYAN)Building Docker images...$(NC)"
	@$(DOCKER_COMPOSE) build

.PHONY: dev-setup
dev-setup: docker-build db-start
	@echo "$(GREEN)Development environment setup complete$(NC)"

help:
	@echo "Available commands:"
	@echo " make build-go    - Build the Go application"
	@echo " make run         - Run the application"
	@echo " make clean       - Clean build files"
	@echo " make test        - Run all tests"
	@echo " make coverage    - Generate test coverage report"
	@echo " make fmt         - Format code"
	@echo " make tidy        - Tidy dependencies"
	@echo " make db-start    - Start PostgreSQL container"
	@echo " make db-stop     - Stop PostgreSQL container"
	@echo " make db-restart  - Restart PostgreSQL container"
	@echo " make db-logs     - View PostgreSQL logs"
	@echo " make db-shell    - Connect to PostgreSQL shell"
	@echo " make dev-setup   - Setup development environment""
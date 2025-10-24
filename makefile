.PHONY: all build run run-hot test test-cover lint tidy clean docker-build help install-tools migration-create migration-up migration-down migration-status

.DEFAULT_GOAL := help

# --- Variables ---
BINARY_NAME := main
TMP_DIR := ./tmp
SERVER_CMD_PATH := ./cmd/server
MIGRATION_CMD_PATH := ./cmd/migration
MIGRATION_DIR := ./internal/app/database/migration

# Tools version
GOOSE_VERSION := latest
AIR_VERSION := latest
GOLANGCI_VERSION := v2.5.0

# Default binary name untuk Linux/macOS
BINARY_FILE := $(TMP_DIR)/$(BINARY_NAME)

# Mendeteksi OS Windows (MINGW atau CYGWIN) dan menambahkan .exe
ifeq ($(findstring MINGW,$(shell uname -s)),MINGW)
    BINARY_FILE := $(TMP_DIR)/$(BINARY_NAME).exe
endif
ifeq ($(findstring CYGWIN,$(shell uname -s)),CYGWIN)
    BINARY_FILE := $(TMP_DIR)/$(BINARY_NAME).exe
endif

# Get GOBIN or GOHOME/bin as instalation path
GOBIN ?= $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN = $(shell go env GOPATH)/bin
endif

# --- Application ---
build: ## Build the production binary
	@echo "Building binary..."
	@mkdir -p build
	@go build -o build/$(BINARY_NAME) $(SERVER_CMD_PATH)/.

run: ## Run the application
	@go run $(SERVER_CMD_PATH)/.

run-hot: ## Run the application with hot-reload (requires 'air')
	@echo "Running with hot-reload (requires 'air' to be installed)..."
	@mkdir -p $(TMP_DIR)
	@air -build.cmd "go build -o $(BINARY_FILE) $(SERVER_CMD_PATH)/." -build.bin "$(BINARY_FILE)"

# --- Dependencies & Tools ---
tidy: ## Tidy go.mod and go.sum
	@echo "Tidying go.mod..."
	@go mod tidy

# --- Tools Installation ---
install-tools: ## Install required Go development tools (goose, air, golangci-lint)
	@echo "Installing development tools to $(GOBIN)..."

	@# Install goose (masih aman via go install)
	@if ! command -v goose &> /dev/null; then \
		echo "Installing goose..."; \
		go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION); \
	else \
		echo "goose is already installed."; \
	fi

	@# Install air (masih aman via go install)
	@if ! command -v air &> /dev/null; then \
		echo "Installing air..."; \
		go install github.com/air-verse/air@$(AIR_VERSION); \
	else \
		echo "air is already installed."; \
	fi

	@# Install golangci-lint (menggunakan metode biner)
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing golangci-lint..."; \
		case "$$(uname -s)" in \
			Darwin) \
				echo "Using 'brew' for macOS..."; \
				brew install golangci-lint; \
				;; \
			Linux | MINGW*) \
				echo "Using 'curl' script for Linux/Windows (MINGW)..."; \
				curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) $(GOLANGCI_VERSION); \
				;; \
			*) \
				echo "Unsupported OS: $$(uname -s) for binary install. Attempting 'go install' as fallback..."; \
				go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION); \
				;; \
		esac; \
	else \
		echo "golangci-lint is already installed."; \
	fi

# --- Quality & Testing ---
test: ## Run all unit tests for internal packages
	@echo "Running tests..."
	@go list ./internal/... | xargs go test -v -cover

test-cover: ## Run all unit tests and open coverage report
	@echo "Running tests with coverage..."
	@go list ./internal/... | xargs go test -coverprofile=coverage.out

	@go tool cover -html=coverage.out

lint: ## Run the linter (golangci-lint)
	@echo "Running linter..."
	@golangci-lint run

# --- Database Migration ---
migration-create: ## Create a new SQL migration file (e.g., make migration-create NAME=create_users_table)
	@if [ -z "$(NAME)" ]; then \
		echo "ERROR: NAME variable is not set."; \
		echo "Usage: make migration-create NAME=create_users_table"; \
		exit 1; \
	fi
	@echo "Creating migration file: $(NAME)"
	@goose -dir $(MIGRATION_DIR) create $(NAME) go

migration-up: ## Run all pending 'up' migrations
	@echo "Running migrations up..."
	@go run $(MIGRATION_CMD_PATH)/. up

migration-down: ## Roll back the last 'up' migration
	@echo "Running migrations down..."
	@go run $(MIGRATION_CMD_PATH)/. down

migration-status: ## Check migration status
	@echo "Checking migration status..."
	@go run $(MIGRATION_CMD_PATH)/. status

migration-fix: ## Apply sequential ordering to migrations
	@echo "Apply sequential ordering to migrations..."
	@go run $(MIGRATION_CMD_PATH)/. fix

# --- Docker & Clean ---
docker-build: build ## Build the production Docker image
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest .

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf build
	@rm -f coverage.out

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

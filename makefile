# This tells Make that these targets aren't actual files.
# This prevents conflicts if a file with the same name exists.
.PHONY: all build run run-hot test test-cover lint tidy clean \
docker-build docker-build-server docker-build-migrator docker-build-all \
podman-build podman-build-server podman-build-migrator podman-build-all \
help install-tools \
migration-create migration-up migration-down migration-status migration-fix

# Sets the default command to run when 'make' is called without arguments.
.DEFAULT_GOAL := help

# --- Variables ---
APP_IMAGE_NAME := travel-api-server
MIGRATOR_IMAGE_NAME := travel-api-migrator
TAG := latest

BINARY_NAME := main
TMP_DIR := ./tmp
SERVER_CMD_PATH := ./cmd/server
MIGRATION_CMD_PATH := ./cmd/migration
MIGRATION_DIR := ./internal/app/database/migration

# Tools version
GOOSE_VERSION := latest
AIR_VERSION := latest
GOLANGCI_VERSION := v2.5.0

# Default binary name for Linux/macOS
BINARY_FILE := $(TMP_DIR)/$(BINARY_NAME)

# Detect Windows OS (MINGW or CYGWIN) and append .exe
ifeq ($(findstring MINGW,$(shell uname -s)),MINGW)
    BINARY_FILE := $(TMP_DIR)/$(BINARY_NAME).exe
endif
ifeq ($(findstring CYGWIN,$(shell uname -s)),CYGWIN)
    BINARY_FILE := $(TMP_DIR)/$(BINARY_NAME).exe
endif

# Get GOBIN (where Go installs binaries), fallback to GOPATH/bin
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

test-cover: ## Run tests and open the HTML coverage report
	@echo "Running tests with coverage..."
	@go list ./internal/... | xargs go test -coverprofile=coverage.out

	@go tool cover -html=coverage.out

lint: ## Run the Go linter (golangci-lint)
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

# --- Docker ---
docker-build: build ## Build the production Container image
	@echo "Building Container image..."
	@docker build -t $(BINARY_NAME):latest .

docker-build-server: ## Build the server container image (production stage)
	@echo "Building server image..."
	@docker build --target server -t $(APP_IMAGE_NAME):$(TAG) .

docker-build-migrator: ## Build the migrator container image (migrator stage)
	@echo "Building migrator image..."
	@docker build --target migrator -t $(MIGRATOR_IMAGE_NAME):$(TAG) .

docker-build-all: docker-build-server docker-build-migrator ## Build all container images (server and migrator)

# --- Podman ---
podman-build: build ## Build the main production container image using Podman
	@echo "Building container image..."
	@podman build -t $(BINARY_NAME):latest .

podman-build-server: ## Build the server container image using Podman
	@echo "Building server image..."
	podman build --target server -t $(APP_IMAGE_NAME):$(TAG) .

podman-build-migrator: ## Build the migrator container image using Podman
	@echo "Building migrator image..."
	@podman build --target migrator -t $(MIGRATOR_IMAGE_NAME):$(TAG) .

podman-build-all: podman-build-server podman-build-migrator ## Build all container images using Podman

clean: ## Clean build artifacts (tmp dir, build dir, coverage)
	@echo "Cleaning build artifacts..."
	@rm -rf build
	@rm -f coverage.out
	@rm -rf tmp

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: help build test clean run install lint fmt vet tidy build-all

# Build variables
BINARY_NAME=edge-client
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build output directory
BUILD_DIR=build

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary for current platform
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/edge-client

build-linux-amd64: ## Build for Linux AMD64
	@echo "Building for Linux AMD64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/edge-client

build-linux-arm64: ## Build for Linux ARM64 (Raspberry Pi)
	@echo "Building for Linux ARM64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/edge-client

build-linux-arm: ## Build for Linux ARM (Raspberry Pi 32-bit)
	@echo "Building for Linux ARM..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm ./cmd/edge-client

build-windows: ## Build for Windows AMD64
	@echo "Building for Windows AMD64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/edge-client

build-macos: ## Build for macOS (both architectures)
	@echo "Building for macOS AMD64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/edge-client
	@echo "Building for macOS ARM64..."
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/edge-client

build-all: build-linux-amd64 build-linux-arm64 build-linux-arm build-windows build-macos ## Build for all platforms
	@echo "All builds complete!"

run: ## Run the application
	$(GOCMD) run ./cmd/edge-client --config configs/config.example.yaml

test: ## Run tests
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

install: build ## Install the binary to $GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
	@echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

fmt: ## Format code
	$(GOFMT) ./...

vet: ## Run go vet
	$(GOVET) ./...

lint: ## Run golangci-lint (requires golangci-lint to be installed)
	golangci-lint run

tidy: ## Tidy go.mod
	$(GOMOD) tidy

deps: ## Download dependencies
	$(GOMOD) download

.DEFAULT_GOAL := help

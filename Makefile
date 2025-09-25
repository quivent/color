# Color CLI Makefile

BINARY_NAME=color
BUILD_DIR=build
INSTALL_DIR=$(HOME)/.local/bin

# Build targets
.PHONY: build clean install uninstall test help

build: ## Build the color CLI
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

install: build ## Install color CLI to ~/.local/bin
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@echo "Installed $(BINARY_NAME) to $(INSTALL_DIR)"
	@echo "Make sure $(INSTALL_DIR) is in your PATH"

uninstall: ## Remove color CLI from ~/.local/bin
	@echo "Removing $(BINARY_NAME) from $(INSTALL_DIR)..."
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Uninstalled $(BINARY_NAME)"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	@echo "Clean complete"

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

fmt: ## Format Go code
	@echo "Formatting Go code..."
	go fmt ./...

lint: ## Run Go linter
	@echo "Running linter..."
	golangci-lint run

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

dev-install: ## Install for development (creates symlink)
	@echo "Installing for development..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@mkdir -p $(INSTALL_DIR)
	ln -sf $(PWD)/$(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Development installation complete"

help: ## Show this help message
	@echo "Color CLI Build System"
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
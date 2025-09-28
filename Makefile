.PHONY: build test dev clean fmt vet lint install help

APP_NAME := tc-server
BUILD_DIR := ./bin
CMD_DIR := ./cmd/tc-server

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

test:
	@echo "Running tests..."
	@go test -v ./...

dev:
	@echo "Running in development mode..."
	@go run $(CMD_DIR)/tc-server.go

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

fmt:
	@echo "Formatting code..."
	@go fmt ./...

vet:
	@echo "Vetting code..."
	@go vet ./...

lint: fmt vet

install:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  test          - Run tests"
	@echo "  dev           - Run in development mode"
	@echo "  clean         - Clean build artifacts"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  lint          - Run fmt and vet"
	@echo "  install       - Install dependencies"
	@echo "  help          - Show this help"
.PHONY: build test lint fmt clean run install

# Binary name
BINARY_NAME=mcp-browser-server
BINARY_PATH=bin/$(BINARY_NAME)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofumpt
GOLINT=golangci-lint

# Build the project
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GOBUILD) -o $(BINARY_PATH) cmd/mcp-browser-server/main.go

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -race -cover ./...

# Run linter
lint:
	@echo "Running linter..."
	$(GOLINT) run

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -w .

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out

# Run the server
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_PATH)

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BINARY_PATH) $(GOPATH)/bin/

# Update dependencies
deps:
	@echo "Updating dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Run all checks
check: fmt lint test

# Help
help:
	@echo "Available commands:"
	@echo "  make build    - Build the project"
	@echo "  make test     - Run tests"
	@echo "  make lint     - Run linter"
	@echo "  make fmt      - Format code"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make run      - Build and run the server"
	@echo "  make install  - Install the binary"
	@echo "  make deps     - Update dependencies"
	@echo "  make check    - Run all checks (fmt, lint, test)"
	@echo "  make help     - Show this help message"
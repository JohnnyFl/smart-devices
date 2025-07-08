# Smart Devices Management System - Makefile

.PHONY: help build test clean deploy dev setup

# Default target
help:
	@echo "Smart Devices Management System"
	@echo ""
	@echo "Available targets:"
	@echo "  build       - Build all Lambda functions"
	@echo "  test        - Run all tests"
	@echo "  test-cover  - Run tests with coverage"
	@echo "  clean       - Clean build artifacts"
	@echo "  deploy      - Deploy to AWS (dev stage)"
	@echo "  deploy-prod - Deploy to AWS (prod stage)"
	@echo "  dev         - Start local development environment"
	@echo "  setup       - Setup local development environment"
	@echo "  lint        - Run Go linter"
	@echo "  fmt         - Format Go code"

# Build all Lambda functions
build:
	@echo "Building Lambda functions..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o bin/get-device cmd/get-device/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o bin/create-device cmd/create-device/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o bin/update-device cmd/update-device/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o bin/delete-device cmd/delete-device/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o bin/list-devices cmd/list-devices/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o bin/sqs-listener cmd/sqs-listener/main.go
	@echo "Build complete!"

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Integration tests
test-integration:
	@echo "Running integration tests..."
	./run_test.sh

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf build/
	rm -rf .serverless/
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Deploy to AWS (development)
deploy:
	@echo "Deploying to AWS (dev stage)..."
	serverless deploy --stage dev

# Deploy to AWS (production)
deploy-prod:
	@echo "Deploying to AWS (prod stage)..."
	serverless deploy --stage prod

# Setup local development environment
setup:
	@echo "Setting up local development environment..."
	@echo "1. Installing dependencies..."
	go mod download
	npm install
	@echo "2. Starting DynamoDB Local..."
	docker-compose up -d
	@echo "3. Waiting for DynamoDB to be ready..."
	sleep 5
	@echo "4. Creating devices table..."
	aws dynamodb create-table \
		--table-name devices \
		--attribute-definitions AttributeName=id,AttributeType=S \
		--key-schema AttributeName=id,KeyType=HASH \
		--billing-mode PAY_PER_REQUEST \
		--endpoint-url http://localhost:8000 || true
	@echo "Setup complete! Run 'make dev' to start development server."

# Start local development
dev:
	@echo "Starting local development server..."
	serverless offline start
# Show project status
status:
	@echo "Project Status:"
	@echo "Go version: $(shell go version)"
	@echo "Node version: $(shell node --version)"
	@echo "Serverless version: $(shell serverless --version)"
	@echo "Docker status: $(shell docker ps --filter name=dynamodb-local --format 'table {{.Names}}\t{{.Status}}')"

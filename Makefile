# Variables
BUILD_DIR := build/
BINARY_NAME := pismo-backend

# Phony targets
.PHONY: all clean deps build run

# Default target
all: clean deps build

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Download dependencies
deps:
	go mod download

# Build the application
build:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux go build -o $(BUILD_DIR)$(BINARY_NAME) ./cmd/

# Run the application inside docker
run:
	docker-compose up -d --build

stop:
	docker-compose down

swagger:
	swag fmt -d ./cmd,./internal
	swag init -g ./cmd/main.go --parseDependency --parseInternal

mocks: ## Generate mocks
	mockgen -destination=internal/service/mock_services/mock.go -package=mockService github.com/akhiltak/pismo-api/internal/service TransactionService
	mockgen -destination=internal/storage/repo/mock_repo/mock.go -package=mockRepo github.com/akhiltak/pismo-api/internal/storage/repo Account,Transaction,Operation

# Test the application
test:
	go test ./... -v

int-test:
	docker-compose -p pismo-test-container -f docker-compose.test.yml up -d --build
	go test ./... -v -tags=integration
	docker-compose -p pismo-test-container -f docker-compose.test.yml down

# Lint the code
lint:
	golangci-lint run

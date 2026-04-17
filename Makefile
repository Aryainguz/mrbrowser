BINARY=mr-browser
BUILD_DIR=./bin
GO=go
GOFLAGS=-ldflags="-s -w"

.PHONY: all build test test-unit test-int test-e2e clean docker lint fmt vet deps tidy

all: build

## Build the mr-browser binary
build:
	@echo "→ Building $(BINARY)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY) ./cli
	@echo "✓ Binary: $(BUILD_DIR)/$(BINARY)"

## Run all tests
test: test-unit

## Run unit tests only
test-unit:
	@echo "→ Running unit tests..."
	$(GO) test -v -race -count=1 ./tests/unit/... ./intelligence/... ./memory/... ./core/...
	@echo "✓ Unit tests passed"

## Run integration tests (requires Chromium)
test-int:
	@echo "→ Running integration tests..."
	$(GO) test -v -timeout 60s ./tests/integration/...
	@echo "✓ Integration tests passed"

## Run e2e tests
test-e2e:
	@echo "→ Running e2e tests..."
	$(GO) test -v -timeout 120s ./tests/e2e/...
	@echo "✓ E2E tests passed"

## Download Go module dependencies
deps:
	$(GO) mod download

## Tidy go.mod and go.sum
tidy:
	$(GO) mod tidy

## Format code
fmt:
	$(GO) fmt ./...

## Run vet
vet:
	$(GO) vet ./...

## Lint (requires golangci-lint)
lint: vet
	@which golangci-lint > /dev/null 2>&1 || (echo "golangci-lint not installed, skipping" && exit 0)
	golangci-lint run ./...

## Build Docker image
docker:
	@echo "→ Building Docker image..."
	docker build -f docker/Dockerfile -t mrbrowser:latest .
	@echo "✓ Docker image: mrbrowser:latest"

## Start with Docker Compose
docker-up:
	docker compose -f docker/docker-compose.yml up -d

## Stop Docker Compose
docker-down:
	docker compose -f docker/docker-compose.yml down

## Clean build artifacts
clean:
	@rm -rf $(BUILD_DIR)
	@echo "✓ Cleaned"

## Install binary to /usr/local/bin
install: build
	cp $(BUILD_DIR)/$(BINARY) /usr/local/bin/$(BINARY)
	@echo "✓ Installed to /usr/local/bin/$(BINARY)"

## Show help
help:
	@grep -E '^## ' Makefile | sed 's/## //'

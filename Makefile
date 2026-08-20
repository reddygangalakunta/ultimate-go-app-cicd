# ==============================================================================
# Enterprise Go Microservice Makefile
# ==============================================================================

APP_NAME ?= ultimate-go-app
VERSION ?= $(shell cat VERSION)
COMMIT_SHA ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
IMAGE_REGISTRY ?= yashwanth1232
IMAGE_TAG ?= $(VERSION)-$(COMMIT_SHA)

.PHONY: all build run test coverage lint docker-build update-tag clean

all: lint test build

## build: Build static application binary
build:
	@echo "--> Building static binary..."
	@mkdir -p bin
	CGO_ENABLED=0 go build -ldflags="-s -w -X 'main.Version=$(VERSION)-$(COMMIT_SHA)'" -o bin/$(APP_NAME) ./cmd/server

## run: Run service locally
run:
	@echo "--> Running service locally..."
	go run ./cmd/server

## test: Run unit tests with race detector
test:
	@echo "--> Running unit tests..."
	go test -v -race ./...

## coverage: Generate code coverage report
coverage:
	@echo "--> Generating coverage report..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report written to coverage.html"

## lint: Run golangci-lint
lint:
	@echo "--> Running golangci-lint..."
	golangci-lint run ./...

## docker-build: Build Docker image
docker-build:
	@echo "--> Building Docker image $(IMAGE_REGISTRY)/$(APP_NAME):$(IMAGE_TAG)..."
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT_SHA=$(COMMIT_SHA) \
		-t $(IMAGE_REGISTRY)/$(APP_NAME):$(IMAGE_TAG) \
		-t $(IMAGE_REGISTRY)/$(APP_NAME):latest .

## update-tag: Update manifest image tag using shell script
update-tag:
	@echo "--> Executing update-image-tag.sh..."
	DRY_RUN=true ./scripts/update-image-tag.sh "$(IMAGE_TAG)" "deployments/k8s/deployment.yaml"

## clean: Clean up build artifacts
clean:
	@echo "--> Cleaning up build artifacts..."
	rm -rf bin/ coverage.out coverage.html

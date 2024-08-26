# Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Docker parameters
DOCKER_COMPOSE=docker-compose
DOCKER=docker

# Kubernetes parameters
KUBECTL=kubectl

# Service names
SERVICES=device user

# Docker image parameters
DOCKER_REGISTRY=your-registry.com
VERSION=$(shell git describe --tags --always --dirty)

.PHONY: all build clean test coverage run deps generate docker-build docker-push deploy $(SERVICES)

all: test build

build: $(SERVICES)

$(SERVICES):
	$(GOBUILD) -o bin/$@ -v ./cmd/$@-service

clean:
	$(GOCLEAN)
	rm -f bin/*

test:
	$(GOTEST) -v ./...

coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

run-%:
	$(GOBUILD) -o bin/$* -v ./cmd/$*-service
	./bin/$*

deps:
	$(GOGET) -v -t -d ./...
	$(GOMOD) tidy

generate:
	$(GOCMD) generate ./...

# Integration tests
integration-test-%:
	$(DOCKER_COMPOSE) -f deployments/docker-compose.$*-test.yml up --build --abort-on-container-exit
	$(DOCKER_COMPOSE) -f deployments/docker-compose.$*-test.yml down

integration-test-all: $(addprefix integration-test-,$(SERVICES))

# Docker commands
docker-build-%:
	$(DOCKER) build -t $(DOCKER_REGISTRY)/$*-service:$(VERSION) -f build/Dockerfile.$* .

docker-push-%: docker-build-%
	$(DOCKER) push $(DOCKER_REGISTRY)/$*-service:$(VERSION)

docker-build-all: $(addprefix docker-build-,$(SERVICES))
docker-push-all: $(addprefix docker-push-,$(SERVICES))

# Kubernetes deployment commands
k8s-deploy-%:
	$(KUBECTL) apply -f deployments/kubernetes/$*-service

k8s-deploy-all: $(addprefix k8s-deploy-,$(SERVICES))

k8s-deploy-shared:
	$(KUBECTL) apply -f deployments/kubernetes/shared

k8s-delete-%:
	$(KUBECTL) delete -f deployments/kubernetes/$*-service

k8s-delete-all: $(addprefix k8s-delete-,$(SERVICES))

k8s-delete-shared:
	$(KUBECTL) delete -f deployments/kubernetes/shared

# Database migrations (example for device service)
migrate-up-device:
	migrate -path internal/device/migrations -database "postgresql://username:password@localhost:5432/device_db?sslmode=disable" up

migrate-down-device:
	migrate -path internal/device/migrations -database "postgresql://username:password@localhost:5432/device_db?sslmode=disable" down

# Linting
lint:
	golangci-lint run

# Proto generation
proto:
	protoc --go_out=. --go-grpc_out=. pkg/proto/*.proto

# Help
help:
	@echo "make - Runs the tests and builds all services"
	@echo "make build - Builds all services"
	@echo "make [service-name] - Builds a specific service"
	@echo "make clean - Removes all binaries and cleans the Go cache"
	@echo "make test - Runs all unit tests"
	@echo "make coverage - Generates a coverage report"
	@echo "make run-[service-name] - Builds and runs a specific service"
	@echo "make deps - Downloads the dependencies"
	@echo "make generate - Runs go generate"
	@echo "make integration-test-[service-name] - Runs integration tests for a specific service"
	@echo "make integration-test-all - Runs integration tests for all services"
	@echo "make docker-build-[service-name] - Builds Docker image for a specific service"
	@echo "make docker-push-[service-name] - Pushes Docker image for a specific service"
	@echo "make docker-build-all - Builds Docker images for all services"
	@echo "make docker-push-all - Pushes Docker images for all services"
	@echo "make k8s-deploy-[service-name] - Deploys a specific service to Kubernetes"
	@echo "make k8s-deploy-all - Deploys all services to Kubernetes"
	@echo "make k8s-deploy-shared - Deploys shared resources to Kubernetes"
	@echo "make k8s-delete-[service-name] - Deletes a specific service from Kubernetes"
	@echo "make k8s-delete-all - Deletes all services from Kubernetes"
	@echo "make k8s-delete-shared - Deletes shared resources from Kubernetes"
	@echo "make migrate-up-[service-name] - Runs database migrations up for a specific service"
	@echo "make migrate-down-[service-name] - Runs database migrations down for a specific service"
	@echo "make lint - Runs the linter"
	@echo "make proto - Generates Go code from proto files"
.PHONY: all
all: generate build

# Generate protobuf code
.PHONY: generate
generate:
	@echo "Generating protobuf code..."
	@for dir in api/*/v1; do \
		if [ -f "$$dir/*.proto" ]; then \
			echo "Generating $$dir..."; \
			cd $$dir && protoc --go_out=. --go_opt=paths=source_relative \
				--go-grpc_out=. --go-grpc_opt=paths=source_relative \
				*.proto; \
			cd - > /dev/null; \
		fi; \
	done

# Build all services
.PHONY: build
build: build-ai build-patient build-doctor build-medical build-department build-auth

.PHONY: build-ai
build-ai:
	@echo "Building AI service..."
	@cd service/ai && go build -o ../../bin/ai-service ./cmd/

.PHONY: build-patient
build-patient:
	@echo "Building Patient service..."
	@cd service/patient && go build -o ../../bin/patient-service ./cmd/

.PHONY: build-doctor
build-doctor:
	@echo "Building Doctor service..."
	@cd service/doctor && go build -o ../../bin/doctor-service ./cmd/

.PHONY: build-medical
build-medical:
	@echo "Building Medical service..."
	@cd service/medical && go build -o ../../bin/medical-service ./cmd/

.PHONY: build-department
build-department:
	@echo "Building Department service..."
	@cd service/department && go build -o ../../bin/department-service ./cmd/

.PHONY: build-auth
build-auth:
	@echo "Building Auth service..."
	@cd service/auth && go build -o ../../bin/auth-service ./cmd/

# Test
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf bin/

# Docker
.PHONY: docker-build
docker-build:
	@echo "Building Docker images..."
	@docker-compose -f deploy/docker-compose.yml build

.PHONY: docker-up
docker-up:
	@echo "Starting services..."
	@docker-compose -f deploy/docker-compose.yml up -d

.PHONY: docker-down
docker-down:
	@echo "Stopping services..."
	@docker-compose -f deploy/docker-compose.yml down

# Lint
.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

# Wire (dependency injection)
.PHONY: wire
wire:
	@echo "Running wire..."
	@cd service/ai && wire ./...
	@cd service/patient && wire ./...
	@cd service/doctor && wire ./...
	@cd service/medical && wire ./...

# Init project structure
.PHONY: init
init:
	@echo "Initializing project..."
	@go mod download
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/google/wire/cmd/wire@latest

# Run services locally
.PHONY: run-ai
run-ai:
	@go run ./service/ai/cmd/ -conf ./service/ai/configs

.PHONY: run-patient
run-patient:
	@go run ./service/patient/cmd/ -conf ./service/patient/configs

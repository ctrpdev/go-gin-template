include .env
export

.PHONY: all build run test clean docker-up docker-down sqlc lint migrate-up migrate-down

# Variables
APP_NAME ?= api
BUILD_DIR ?= build
MAIN_FILE ?= ./cmd/api/main.go
DB_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

all: build

# Compilar
build:
	@echo "Building..."
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

# Ejecutar local 
run: build
	@echo "Running..."
	./$(BUILD_DIR)/$(APP_NAME)

# Testear
test:
	@echo "Testing..."
	go test -v ./...

# Limpiar compilados
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)

# Docker
docker-dev:
	docker-compose -f deployment/docker-compose.yml -f deployment/docker-compose.dev.yml up -d --build

docker-prod:
	docker-compose -f deployment/docker-compose.yml up -d --build

docker-down:
	docker-compose down

# Sqlc
sqlc:
	sqlc generate

# Go mod tidy and fmt
lint:
	go mod tidy
	go fmt ./...

# Validating and Running Migrations
migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" -verbose down

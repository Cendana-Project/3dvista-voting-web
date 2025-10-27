.PHONY: help run build test clean migrate-up migrate-down seed docker-build docker-up docker-down docker-logs

# Default target
help:
	@echo "Available targets:"
	@echo "  run           - Run the application locally"
	@echo "  build         - Build the application binary"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  seed          - Seed database with initial data"
	@echo "  docker-build  - Build Docker images"
	@echo "  docker-up     - Start services with Docker Compose"
	@echo "  docker-down   - Stop services with Docker Compose"
	@echo "  docker-logs   - View Docker Compose logs"

# Run the application locally
run:
	@echo "Running application..."
	go run ./cmd/server/main.go

# Build the application
build:
	@echo "Building application..."
	go build -o bin/voteweb ./cmd/server

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out

# Run database migrations up
migrate-up:
	@echo "Running migrations..."
	@if [ ! -f ".env" ]; then \
		echo "Error: .env file not found. Please copy env.example to .env and configure it."; \
		exit 1; \
	fi
	@. ./.env && psql $${DATABASE_URL} -f migrations/0001_init.sql

# Rollback database migrations (drop tables)
migrate-down:
	@echo "Rolling back migrations..."
	@if [ ! -f ".env" ]; then \
		echo "Error: .env file not found."; \
		exit 1; \
	fi
	@. ./.env && psql $${DATABASE_URL} -c "DROP TABLE IF EXISTS votes CASCADE; DROP TABLE IF EXISTS innovations CASCADE;"

# Seed database
seed:
	@echo "Seeding database..."
	SEED=true go run ./cmd/server/main.go &
	@sleep 2
	@pkill -f "go run ./cmd/server/main.go" || true

# Docker commands
docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-up:
	@echo "Starting services..."
	SEED=true docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo "Running migrations..."
	docker-compose exec -T db psql -U postgres -d voteweb < migrations/0001_init.sql
	@echo "Services started successfully!"
	@echo "Application available at http://localhost:8080"

docker-down:
	@echo "Stopping services..."
	docker-compose down

docker-logs:
	docker-compose logs -f app

# Development workflow
dev: docker-down docker-up docker-logs



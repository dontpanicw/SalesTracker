.PHONY: run build docker-up docker-down test clean

run:
	go run cmd/main.go

build:
	go build -o bin/analytics cmd/main.go

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down -v

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-integration:
	@echo "Running integration tests..."
	@echo "Starting PostgreSQL container..."
	@docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 3
	@docker exec -it analytics-postgres createdb -U postgres analytics_test 2>/dev/null || true
	@echo "Running tests..."
	@TEST_DATABASE_DSN="host=localhost port=5432 user=postgres password=postgres dbname=analytics_test sslmode=disable" go test -v ./internal/adapter/repository/postgres/
	@echo "Cleaning up..."
	@docker exec -it analytics-postgres dropdb -U postgres analytics_test 2>/dev/null || true

clean:
	rm -rf bin/
	docker-compose down -v

deps:
	go mod download
	go mod tidy

migrate:
	go run cmd/main.go migrate

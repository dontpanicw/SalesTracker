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
	@if [ -z "$$TEST_DATABASE_DSN" ]; then \
		echo "TEST_DATABASE_DSN not set. Using default..."; \
		export TEST_DATABASE_DSN="host=localhost port=5432 user=postgres password=postgres dbname=analytics_test sslmode=disable"; \
	fi
	go test -v ./internal/adapter/repository/postgres/

clean:
	rm -rf bin/
	docker-compose down -v

deps:
	go mod download
	go mod tidy

migrate:
	go run cmd/main.go migrate

.PHONY: setup run clean test lint db-init db-seed db-reset db-migrate db-status help

# Default target
.DEFAULT_GOAL := help

# Variables
GO = go

# Help target
help:
	@echo "Mathtermind (Go Version) Makefile"
	@echo "==============================="
	@echo "Available commands:"
	@echo "  make setup      - Set up the project (install dependencies)"
	@echo "  make run        - Run the application"
	@echo "  make clean      - Clean up temporary files and caches"
	@echo "  make test       - Run all tests"
	@echo "  make lint       - Run linting tools"
	@echo "  make db-init    - Initialize the database"
	@echo "  make db-seed    - Seed the database with sample data"
	@echo "  make db-reset   - Reset the database"
	@echo "  make db-migrate - Run database migrations"
	@echo "  make db-status  - Show database migration status"

# Setup target
setup:
	$(GO) mod tidy

# Run target
run:
	$(GO) run main.go

# Clean target
clean:
	$(GO) clean
	find . -name "*.test" -type f -delete
	find . -name "*.test.exe" -type f -delete
	find . -name "*.test.dSYM" -type d -delete
	find . -name "*.prof" -type f -delete
	find . -name "*.out" -type f -delete
	find . -name "*.cov" -type f -delete

# Test target
test:
	$(GO) test ./...

# Lint target
lint:
	golangci-lint run

# Coverage target
cov:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Database targets
db-init:
	@echo "Database initialization not implemented yet"

db-seed:
	@echo "Database seeding not implemented yet"

db-reset:
	@echo "Database reset not implemented yet"

db-migrate:
	@echo "Database migrations not implemented yet"

db-status:
	@echo "Database status not implemented yet"

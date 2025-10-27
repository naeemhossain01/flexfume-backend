.PHONY: help run build test clean install dev migrate-up migrate-down migrate-rollback migrate-create seed docker-build docker-run docker-stop deploy-render

help: ## Show this help message
	@echo "Available commands:"
	@echo "  make run              - Run the application"
	@echo "  make build            - Build the application"
	@echo "  make test             - Run tests"
	@echo "  make install          - Install dependencies"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make dev              - Run in development mode with auto-reload"
	@echo "  make migrate-up       - Run database migrations"
	@echo "  make migrate-down     - Rollback all migrations"
	@echo "  make migrate-rollback - Rollback last migration"
	@echo "  make migrate-create   - Create a new migration file"
	@echo "  make seed             - Seed database with initial data (admin user)"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make docker-run       - Run Docker container"
	@echo "  make docker-stop      - Stop Docker container"
	@echo "  make deploy-render    - Deploy to Render (requires git push)"

run: ## Run the application
	go run cmd/api/main.go

build: ## Build the application
	go build -o bin/server cmd/api/main.go

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

install: ## Install dependencies
	go mod download
	go mod tidy

dev: ## Run in development mode
	go run cmd/api/main.go

migrate-up: ## Run database migrations
	go run cmd/migrate/main.go -cmd=up

migrate-down: ## Rollback all migrations
	go run cmd/migrate/main.go -cmd=down

migrate-rollback: ## Rollback last migration
	go run cmd/migrate/main.go -cmd=rollback

migrate-create: ## Create a new migration file (usage: make migrate-create name=your_migration_name)
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a migration name using 'make migrate-create name=your_migration_name'"; \
		exit 1; \
	fi
	@timestamp=$$(date +%s); \
	touch migrations/$${timestamp}_$(name).up.sql; \
	touch migrations/$${timestamp}_$(name).down.sql; \
	echo "Created migrations/$${timestamp}_$(name).up.sql"; \
	echo "Created migrations/$${timestamp}_$(name).down.sql"

seed: ## Seed database with initial data
	go run cmd/seed/main.go

docker-build: ## Build Docker image
	docker build -t flexfume-ecom-backend:latest .

docker-run: ## Run Docker container (requires .env file)
	docker run -p 8080:8080 --env-file .env flexfume-ecom-backend:latest

docker-stop: ## Stop Docker container
	docker stop $$(docker ps -q --filter "ancestor=flexfume-ecom-backend:latest") || true

deploy-render: ## Deploy to Render (push to main branch)
	@echo "Render deployment is automatic on git push to main branch"
	@echo "1. Ensure all changes are committed"
	@echo "2. Run: git push origin main"
	@echo "3. Monitor deployment at: https://dashboard.render.com"
	@echo ""
	@echo "For manual deployment:"
	@echo "1. Go to Render Dashboard"
	@echo "2. Select flexfume-ecom-backend service"
	@echo "3. Click 'Manual Deploy' → 'Deploy latest commit'"

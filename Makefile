include .env


.PHONY: protoc
protoc:
	@if not exist pkg\protos\gen\go mkdir pkg\protos\gen\go
	protoc --proto_path=pkg/protos/proto --go_out=pkg/protos/gen/go --go_opt=paths=source_relative --go-grpc_out=pkg/protos/gen/go --go-grpc_opt=paths=source_relative pkg/protos/proto/device/device.proto pkg/protos/proto/payment/payment.proto pkg/protos/proto/refund/refund.proto pkg/protos/proto/refund_enhanced/refund_enhanced.proto pkg/protos/proto/utility/utility.proto


.PHONY: db/migrations
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=sql -dir=./migrations ${name}

db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path=./migrations -database=${database} up

db/migrations/down:
	@echo 'Running up migrations ...'
	migrate -path=./migrations -database=${POSTGRES_URI} down 1

.PHONY: db/dump
db/dump/create:
	docker compose run --rm db-tools /db_dump.sh

db/dump/restore:
	docker compose run --rm db-tools /db_restore.sh

POSTGRES_URI = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)
setup:
	@echo "🚀 Setting up the project..."
	@if [ ! -f .env ]; then \
		echo "Creating .env file from example..."; \
		cp .env.example .env; \
		echo ".env file created. You may want to edit it before continuing."; \
		echo "Press Enter to continue or Ctrl+C to abort and edit .env first."; \
		read dummy; \
		include .env; \
		export; \
	fi
	@echo "Creating dumps directory..."
	@mkdir -p dumps
	@echo "Building Docker images..."
	docker-compose build
	@echo "Starting database..."
	docker-compose up -d db
	@echo "Waiting for database to initialize..."
	@sleep 5
	@if [ -f ./dumps/kaspi_pay.custom ]; then \
		echo "Found existing database dump, restoring..."; \
		$(MAKE) db/dump/restore; \
	else \
		echo "No existing dump found, creating initial database..."; \
		echo "Running migrations..."; \
		DB_CONTAINER_HOST=$$(docker-compose port db 5432 | cut -d ':' -f 1); \
		CONTAINER_POSTGRES_URI="postgres://$(DB_USER):$(DB_PASSWORD)@$$DB_CONTAINER_HOST:$$(docker-compose port db 5432 | cut -d ':' -f 2)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)"; \
		echo "Using database connection: $$CONTAINER_POSTGRES_URI"; \
		migrate -path=./migrations -database=$$CONTAINER_POSTGRES_URI up; \
		echo "Creating initial dump..."; \
		$(MAKE) db/dump/create; \
	fi
	@echo "Starting all services..."
	docker-compose up -d
	@echo "	Setup complete! The application is now running."
	@echo "   API is available at http://localhost:8080"
	@echo "   You can view logs with: docker-compose logs -f"
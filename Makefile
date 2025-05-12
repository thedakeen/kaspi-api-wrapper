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
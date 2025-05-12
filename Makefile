.PHONY: db/migrations
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=sql -dir=./migrations ${name}

db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path=./migrations -database=${POSTGRES_URI} up

db/migrations/down:
	@echo 'Running up migrations ...'
	migrate -path=./migrations -database=${POSTGRES_URI} down 1
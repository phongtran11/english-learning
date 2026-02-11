include .env
export

MIGRATE_DSN ?= "user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) host=$(DB_HOST) port=$(DB_PORT) sslmode=require"

.PHONY: migrate-up migrate-down migrate-create run build watch

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

watch:
	air

migrate-up:
	goose -dir migrations postgres $(MIGRATE_DSN) up

migrate-down:
	goose -dir migrations postgres $(MIGRATE_DSN) down

migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir migrations create $$name sql

migrate-status:
	goose -dir migrations postgres $(MIGRATE_DSN) status

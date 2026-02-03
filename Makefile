-include .env
export $(shell [ -f .env ] && sed 's/=.*//' .env)

MIGRATIONS_DIR=./migrations
DATABASE_URL=$(POSTGRES_DSN)

.PHONY: help migrate-create migrate-up migrate-down migrate-force migrate-version

help:
	@echo "Available commands:"
	@echo "  make migrate-create name=<name>  Create a new migration file"
	@echo "  make migrate-up                  Run all up migrations"
	@echo "  make migrate-down                Rollback the last migration"
	@echo "  make migrate-force version=<v>   Force a specific version"
	@echo "  make migrate-version             Check current migration version"

migrate-create:
	@if [ -z "$(name)" ]; then echo "Error: name is required. Usage: make migrate-create name=<name>"; exit 1; fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	@echo "Running up migrations..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" -verbose up

migrate-down:
	@echo "Running down migrations..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" -verbose down 1

migrate-force:
	@if [ -z "$(version)" ]; then echo "Error: version is required. Usage: make migrate-force version=<v>"; exit 1; fi
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" -verbose force $(version)

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

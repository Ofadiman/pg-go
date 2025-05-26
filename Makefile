# Makes all possible targets in the Makefile non-file targets, ensuring they are always executed regardless of the presence of files with matching names.
.PHONY: *

# Set environment variable for all shell commands run in the scope of makefile. This variable is not set in the main shell. 
export GOOSE_DBSTRING = postgres://admin:password@localhost:5432/pg-go
export GOOSE_DRIVER = postgres
export GOOSE_MIGRATION_DIR = ./migrations

database_recreate_public_schema:
	docker compose exec postgres psql -U admin -d pg-go -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

database_migrate:
	goose up

database_seed:
	go run src/seed/main.go

# make database_new_migration name=create_users_table
database_new_migration:
	goose create $(name) sql

database_recreate: database_recreate_public_schema database_migrate database_seed

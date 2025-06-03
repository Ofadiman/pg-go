# This sets the `-s` flag for every make invocation, ensuring that all targets run silently.
MAKEFLAGS += -s

# Makes all possible targets in the Makefile non-file targets, ensuring they are always executed regardless of the presence of files with matching names.
.PHONY: *

$(shell command -v goose >/dev/null 2>&1 || { echo "The goose package is not installed. Please install it using the official documentation: https://github.com/pressly/goose" >&2; exit 1; })

# Set environment variable for all shell commands run in the scope of makefile. This variable is not set in the main shell. 
export GOOSE_DBSTRING = postgres://admin:password@localhost:5432/pg-go
export GOOSE_DRIVER = postgres
export GOOSE_MIGRATION_DIR = ./migrations

help:  ## Show this help menu.
	grep -E '^[a-zA-Z0-9_-]+:.*##' $(MAKEFILE_LIST) | \
	sort | \
	awk -F ':.*##' ' \
		NR==1{max=0} \
		{if(length($$1)>max) max=length($$1); lines[NR]=$$0} \
		END { \
			for(i=1; i<=NR; i++) { \
				split(lines[i], parts, ":.*##"); \
				printf "\033[36m%-*s\033[0m %s\n", max+2, parts[1], parts[2] \
			} \
		}'

database_recreate_public_schema:  ## Drop and recreate the public schema in the database.
	docker compose exec postgres psql -U admin -d pg-go -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

database_migrate:  ## Run all pending database migrations.
	goose up

database_seed:  ## Seed the database with initial data.
	go run cmd/seed/main.go

database_new_migration:  ## Create a new migration. Usage: "make database_new_migration name=your_migration_name".
	goose create $(name) sql

database_recreate: database_recreate_public_schema database_migrate database_seed  ## Recreate public schema, run migrations, and seed the database.

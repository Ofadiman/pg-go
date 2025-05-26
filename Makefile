.PHONY: * help

# Set environment variable for all shell commands run in the scope of makefile. This variable is not set in the main shell. 
export GOOSE_DBSTRING = postgres://admin:password@localhost:5432/pg-go
export GOOSE_DRIVER = postgres
export GOOSE_MIGRATION_DIR = ./migrations

help:  ## Show this help menu.
	@awk 'BEGIN {FS = ":.*##"} \
		/^[a-zA-Z0-9_-]+:.*##/ { \
			if(length($$1)>max) max=length($$1); \
			help[$$1]=$$2; \
			targets[++n]=$$1 \
		} \
		END { \
			asort(targets); \
			for (i=1; i<=n; i++) \
				printf "\033[36m%-*s\033[0m %s\n", max+2, targets[i], help[targets[i]] \
		}' $(MAKEFILE_LIST)

database_recreate_public_schema:  ## Drop and recreate the public schema in the database.
	docker compose exec postgres psql -U admin -d pg-go -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

database_migrate:  ## Run all pending database migrations.
	goose up

database_seed:  ## Seed the database with initial data.
	go run src/seed/main.go

database_new_migration:  ## Create a new migration. Usage: "make database_new_migration name=your_migration_name".
	goose create $(name) sql

database_recreate: database_recreate_public_schema database_migrate database_seed  ## Recreate public schema, run migrations, and seed the database.

database_reset_schema:
	docker compose exec postgres psql -U admin -d pg-go -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

database_migrate:
	goose up

database_seed:
	go run src/seed/main.go

database_recreate: database_reset_schema database_migrate database_seed

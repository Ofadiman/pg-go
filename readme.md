# pg-go

**pg-go** is a benchmarking project that compares various query builder libraries available in Go (Golang).

## Database Migrations

This project uses [Goose](https://github.com/pressly/goose) for managing database migrations.

### Creating a New Migration

To create a new migration file, navigate to the `migrations` directory and run:

```sh
goose create create_users_table sql
```

This command will generate a new SQL migration file with a timestamped filename.

### Applying Migrations

To apply all pending migrations, execute the following command from the project root:

```sh
goose up
```

This will run all migrations in the correct order.

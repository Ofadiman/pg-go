package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ofadiman/pg-go/migrations"
	"github.com/ofadiman/pg-go/seed"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const image = "postgres:17.4-bookworm@sha256:7f29c02ba9eeff4de9a9f414d803faa0e6fe5e8d15ebe217e3e418c82e652b35"

var now = time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

var (
	container *postgres.PostgresContainer
	dsn       string
)

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	time.Local = time.UTC
	ctx := context.Background()

	var err error
	container, err = postgres.Run(ctx, image,
		postgres.WithDatabase("pg-go"),
		postgres.WithUsername("admin"),
		postgres.WithPassword("password"),
		postgres.WithSQLDriver("pgx"),
		postgres.BasicWaitStrategies(),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			fmt.Fprintln(os.Stderr, "terminate container:", err)
		}
	}()
	if err != nil {
		fmt.Fprintln(os.Stderr, "start container:", err)
		return 1
	}

	dsn, err = container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintln(os.Stderr, "connection string:", err)
		return 1
	}

	if err := prepare(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "prepare database:", err)
		return 1
	}

	if err := container.Snapshot(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "snapshot database:", err)
		return 1
	}

	return m.Run()
}

func prepare(ctx context.Context) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrations.FS)
	if err != nil {
		return err
	}
	if _, err := provider.Up(ctx); err != nil {
		return err
	}

	return seed.Run(ctx, db)
}

func newDB(t *testing.T) *sql.DB {
	t.Helper()

	if err := container.Restore(t.Context()); err != nil {
		t.Fatalf("restore database: %v", err)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return db
}

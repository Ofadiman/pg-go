package gormrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ofadiman/pg-go/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Store struct {
	q   *gorm.DB
	now func() time.Time
}

func New(db *sql.DB, now func() time.Time) (repository.Store, error) {
	database, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open gorm: %w", err)
	}
	return &Store{q: database, now: now}, nil
}

func (s *Store) Users() repository.UserRepository {
	return users{s}
}

func (s *Store) Posts() repository.PostRepository {
	return posts{s}
}

func (s *Store) Comments() repository.CommentRepository {
	return comments{s}
}

func (s *Store) InTx(ctx context.Context, fn func(tx repository.Store) error) error {
	return s.q.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&Store{q: tx, now: s.now})
	})
}

func notFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func translate(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "unique_email" {
		return fmt.Errorf("%w: %w", repository.ErrEmailTaken, err)
	}
	return err
}

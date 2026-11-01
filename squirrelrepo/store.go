package squirrelrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/ofadiman/pg-go/repository"
)

type Store struct {
	db         *sqlx.DB
	q          sqlx.ExtContext
	tx         *sqlx.Tx
	sb         squirrel.StatementBuilderType
	now        func() time.Time
	savepoints int
}

func New(db *sql.DB, now func() time.Time) (repository.Store, error) {
	database := sqlx.NewDb(db, "pgx")
	return &Store{
		db:  database,
		q:   database,
		sb:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		now: now,
	}, nil
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
	if s.tx != nil {
		return s.inSavepoint(ctx, fn)
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(&Store{db: s.db, q: tx, tx: tx, sb: s.sb, now: s.now}); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(err, fmt.Errorf("rollback transaction: %w", rollbackErr))
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func (s *Store) inSavepoint(ctx context.Context, fn func(tx repository.Store) error) error {
	s.savepoints++
	name := fmt.Sprintf("sp_%d", s.savepoints)

	if _, err := s.tx.ExecContext(ctx, "SAVEPOINT "+name); err != nil {
		return fmt.Errorf("create savepoint: %w", err)
	}

	if err := fn(s); err != nil {
		if _, rollbackErr := s.tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT "+name); rollbackErr != nil {
			return errors.Join(err, fmt.Errorf("rollback to savepoint: %w", rollbackErr))
		}
		return err
	}

	if _, err := s.tx.ExecContext(ctx, "RELEASE SAVEPOINT "+name); err != nil {
		return fmt.Errorf("release savepoint: %w", err)
	}
	return nil
}

func (s *Store) get(ctx context.Context, dest any, query squirrel.Sqlizer) error {
	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}
	return sqlx.GetContext(ctx, s.q, dest, sql, args...)
}

func (s *Store) selectAll(ctx context.Context, dest any, query squirrel.Sqlizer) error {
	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}
	return sqlx.SelectContext(ctx, s.q, dest, sql, args...)
}

func notFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func translate(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "unique_email" {
		return fmt.Errorf("%w: %w", repository.ErrEmailTaken, err)
	}
	return err
}

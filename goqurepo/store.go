package goqurepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ofadiman/pg-go/repository"
)

type queryer interface {
	From(from ...any) *goqu.SelectDataset
	Insert(table any) *goqu.InsertDataset
	Update(table any) *goqu.UpdateDataset
	Delete(table any) *goqu.DeleteDataset
}

type Store struct {
	db         *goqu.Database
	q          queryer
	tx         *goqu.TxDatabase
	now        func() time.Time
	savepoints int
}

func New(db *sql.DB, now func() time.Time) (repository.Store, error) {
	database := goqu.New("postgres", db)
	return &Store{db: database, q: database, now: now}, nil
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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(&Store{db: s.db, q: tx, tx: tx, now: s.now}); err != nil {
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

func translate(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "unique_email" {
		return fmt.Errorf("%w: %w", repository.ErrEmailTaken, err)
	}
	return err
}

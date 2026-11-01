package integration

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ofadiman/pg-go/repository"
)

func TestTransactions(t *testing.T) {
	scenario(t, "tx_commit", func(t *testing.T, r *run) {
		type result struct {
			User  repository.User
			Posts []repository.Post
		}
		got, err := record(r, func() (result, error) {
			var res result
			err := r.store.InTx(t.Context(), func(tx repository.Store) error {
				user, err := tx.Users().Create(t.Context(), newUser(1))
				if err != nil {
					return err
				}
				posts, err := tx.Posts().CreateMany(t.Context(), []repository.NewPost{
					newPost(user.ID, 1),
					newPost(user.ID, 2),
					newPost(user.ID, 3),
				})
				if err != nil {
					return err
				}
				res = result{User: user, Posts: posts}
				return nil
			})
			return res, err
		})
		if err != nil {
			t.Fatal(err)
		}
		r.markGenerated(t, got.User.ID)
		for _, post := range got.Posts {
			r.markGenerated(t, post.ID)
		}

		if _, err := r.store.Users().GetByID(t.Context(), got.User.ID); err != nil {
			t.Errorf("committed user is not readable: %v", err)
		}
		checkResult(t, r, got)
	})

	scenario(t, "tx_rollback", func(t *testing.T, r *run) {
		var userID string
		err := r.record(func() error {
			return r.store.InTx(t.Context(), func(tx repository.Store) error {
				user, err := tx.Users().Create(t.Context(), newUser(1))
				if err != nil {
					return err
				}
				userID = user.ID
				_, err = tx.Posts().CreateMany(t.Context(), []repository.NewPost{
					newPost(user.ID, 1),
					newPost(missingUserID, 2),
				})
				return err
			})
		})

		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || pgErr.Code != "23503" {
			t.Fatalf("got error %v, want a foreign key violation", err)
		}
		r.markGenerated(t, userID)

		if _, err := r.store.Users().GetByID(t.Context(), userID); !errors.Is(err, repository.ErrNotFound) {
			t.Errorf("user created in a rolled back transaction is readable: %v", err)
		}
	})

	scenario(t, "tx_savepoint", func(t *testing.T, r *run) {
		errInner := errors.New("inner transaction failed")
		var outerID, innerID string
		err := r.record(func() error {
			return r.store.InTx(t.Context(), func(tx repository.Store) error {
				outer, err := tx.Users().Create(t.Context(), newUser(1))
				if err != nil {
					return err
				}
				outerID = outer.ID

				err = tx.InTx(t.Context(), func(tx repository.Store) error {
					inner, err := tx.Users().Create(t.Context(), newUser(2))
					if err != nil {
						return err
					}
					innerID = inner.ID
					return errInner
				})
				if !errors.Is(err, errInner) {
					return fmt.Errorf("inner transaction: got error %v, want %v", err, errInner)
				}
				return nil
			})
		})
		if err != nil {
			t.Fatal(err)
		}
		r.markGenerated(t, outerID, innerID)

		if _, err := r.store.Users().GetByID(t.Context(), outerID); err != nil {
			t.Errorf("user created in the outer transaction is not readable: %v", err)
		}
		if _, err := r.store.Users().GetByID(t.Context(), innerID); !errors.Is(err, repository.ErrNotFound) {
			t.Errorf("user created in the rolled back savepoint is readable: %v", err)
		}
	})

	scenario(t, "tx_skip_locked", func(t *testing.T, r *run) {
		type result struct {
			First  []repository.User
			Second []repository.User
		}
		got, err := record(r, func() (result, error) {
			var res result
			locked := make(chan struct{})
			release := make(chan struct{})
			done := make(chan error, 1)

			go func() {
				done <- r.store.InTx(t.Context(), func(tx repository.Store) error {
					users, err := tx.Users().ClaimUnverified(t.Context(), 2)
					if err != nil {
						return err
					}
					res.First = users
					close(locked)
					<-release
					return nil
				})
			}()

			select {
			case <-locked:
			case err := <-done:
				return res, err
			}

			err := r.store.InTx(t.Context(), func(tx repository.Store) error {
				users, err := tx.Users().ClaimUnverified(t.Context(), 2)
				res.Second = users
				return err
			})
			close(release)
			return res, errors.Join(err, <-done)
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(got.First) != 2 || len(got.Second) != 2 {
			t.Fatalf("got %d and %d claimed users, want 2 and 2", len(got.First), len(got.Second))
		}
		for _, first := range got.First {
			for _, second := range got.Second {
				if first.ID == second.ID {
					t.Errorf("user %s was claimed by both transactions", first.ID)
				}
			}
		}
		checkResult(t, r, got)
	})
}

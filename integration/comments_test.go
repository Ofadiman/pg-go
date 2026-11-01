package integration

import (
	"errors"
	"testing"

	"github.com/ofadiman/pg-go/repository"
	"github.com/ofadiman/pg-go/seed"
)

func TestComments(t *testing.T) {
	scenario(t, "comments_thread", func(t *testing.T, r *run) {
		got, err := record(r, func() ([]repository.ThreadComment, error) {
			return r.store.Comments().Thread(t.Context(), seed.ThreadRootCommentID)
		})
		if err != nil {
			t.Fatal(err)
		}
		checkResult(t, r, got)
	})

	scenario(t, "comments_thread_not_found", func(t *testing.T, r *run) {
		_, err := record(r, func() ([]repository.ThreadComment, error) {
			return r.store.Comments().Thread(t.Context(), missingCommentID)
		})
		if !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("got error %v, want %v", err, repository.ErrNotFound)
		}
	})
}

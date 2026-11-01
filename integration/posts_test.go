package integration

import (
	"errors"
	"testing"

	"github.com/ofadiman/pg-go/repository"
	"github.com/ofadiman/pg-go/seed"
)

func TestPosts(t *testing.T) {
	scenario(t, "posts_get_with_author", func(t *testing.T, r *run) {
		got, err := record(r, func() (repository.PostWithAuthor, error) {
			return r.store.Posts().GetWithAuthor(t.Context(), seed.ThreadPostID)
		})
		if err != nil {
			t.Fatal(err)
		}
		checkResult(t, r, got)
	})

	scenario(t, "posts_get_with_author_not_found", func(t *testing.T, r *run) {
		_, err := record(r, func() (repository.PostWithAuthor, error) {
			return r.store.Posts().GetWithAuthor(t.Context(), missingPostID)
		})
		if !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("got error %v, want %v", err, repository.ErrNotFound)
		}
	})

	scenario(t, "posts_list_recent_with_comments", func(t *testing.T, r *run) {
		got, err := record(r, func() ([]repository.PostWithComments, error) {
			return r.store.Posts().ListRecentWithComments(t.Context(), 10)
		})
		if err != nil {
			t.Fatal(err)
		}
		checkResult(t, r, got)
	})
}

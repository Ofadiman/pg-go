package integration

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ofadiman/pg-go/repository"
	"github.com/ofadiman/pg-go/seed"
)

func TestUsers(t *testing.T) {
	scenario(t, "users_create", func(t *testing.T, r *run) {
		got, err := record(r, func() (repository.User, error) {
			return r.store.Users().Create(t.Context(), newUser(1))
		})
		if err != nil {
			t.Fatal(err)
		}
		r.markGenerated(t, got.ID)
		checkResult(t, r, got)
	})

	scenario(t, "users_create_email_taken", func(t *testing.T, r *run) {
		in := newUser(1)
		in.Email = "user-1@example.com"
		_, err := record(r, func() (repository.User, error) {
			return r.store.Users().Create(t.Context(), in)
		})
		if !errors.Is(err, repository.ErrEmailTaken) {
			t.Fatalf("got error %v, want %v", err, repository.ErrEmailTaken)
		}
	})

	scenario(t, "users_create_many", func(t *testing.T, r *run) {
		in := make([]repository.NewUser, 100)
		for i := range in {
			in[i] = newUser(i + 1)
		}
		got, err := record(r, func() ([]repository.User, error) {
			return r.store.Users().CreateMany(t.Context(), in)
		})
		if err != nil {
			t.Fatal(err)
		}
		for _, user := range got {
			r.markGenerated(t, user.ID)
		}
		checkResult(t, r, got)
	})

	scenario(t, "users_get_by_id", func(t *testing.T, r *run) {
		got, err := record(r, func() (repository.User, error) {
			return r.store.Users().GetByID(t.Context(), seed.FirstUserID)
		})
		if err != nil {
			t.Fatal(err)
		}
		checkResult(t, r, got)
	})

	scenario(t, "users_get_by_id_not_found", func(t *testing.T, r *run) {
		_, err := record(r, func() (repository.User, error) {
			return r.store.Users().GetByID(t.Context(), missingUserID)
		})
		if !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("got error %v, want %v", err, repository.ErrNotFound)
		}
	})

	scenario(t, "users_update", func(t *testing.T, r *run) {
		got, err := record(r, func() (repository.User, error) {
			return r.store.Users().Update(t.Context(), seed.FirstUserID, repository.UserUpdate{
				LastName:      new("Sparrow"),
				EmailVerified: new(false),
			})
		})
		if err != nil {
			t.Fatal(err)
		}
		checkResult(t, r, got)
	})

	scenario(t, "users_update_not_found", func(t *testing.T, r *run) {
		_, err := record(r, func() (repository.User, error) {
			return r.store.Users().Update(t.Context(), missingUserID, repository.UserUpdate{LastName: new("Sparrow")})
		})
		if !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("got error %v, want %v", err, repository.ErrNotFound)
		}
	})

	scenario(t, "users_delete", func(t *testing.T, r *run) {
		got, err := record(r, func() (repository.User, error) {
			return r.store.Users().Delete(t.Context(), seed.UserWithoutActivityID)
		})
		if err != nil {
			t.Fatal(err)
		}
		checkResult(t, r, got)

		if _, err := r.store.Users().GetByID(t.Context(), seed.UserWithoutActivityID); !errors.Is(err, repository.ErrNotFound) {
			t.Errorf("deleted user is still readable: %v", err)
		}
	})

	scenario(t, "users_delete_not_found", func(t *testing.T, r *run) {
		_, err := record(r, func() (repository.User, error) {
			return r.store.Users().Delete(t.Context(), missingUserID)
		})
		if !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("got error %v, want %v", err, repository.ErrNotFound)
		}
	})

	scenario(t, "users_upsert_insert", func(t *testing.T, r *run) {
		got, err := record(r, func() (repository.User, error) {
			return r.store.Users().Upsert(t.Context(), newUser(1))
		})
		if err != nil {
			t.Fatal(err)
		}
		r.markGenerated(t, got.ID)
		checkResult(t, r, got)
	})

	scenario(t, "users_upsert_update", func(t *testing.T, r *run) {
		in := newUser(2)
		in.Email = "user-1@example.com"
		got, err := record(r, func() (repository.User, error) {
			return r.store.Users().Upsert(t.Context(), in)
		})
		if err != nil {
			t.Fatal(err)
		}
		if got.ID != seed.FirstUserID {
			t.Errorf("got id %s, want %s", got.ID, seed.FirstUserID)
		}
		checkResult(t, r, got)
	})

	listCases := []struct {
		scenario string
		params   repository.ListUsersParams
	}{
		{
			scenario: "users_list_defaults",
			params:   repository.ListUsersParams{},
		},
		{
			scenario: "users_list_filtered",
			params: repository.ListUsersParams{
				FirstNamePrefix: new("An"),
				Gender:          new("female"),
				EmailVerified:   new(true),
				BornAfter:       new(time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)),
				BornBefore:      new(time.Date(2005, time.January, 1, 0, 0, 0, 0, time.UTC)),
				Sort: []repository.UserSort{
					{Field: repository.SortByLastName},
					{Field: repository.SortByBirthday, Desc: true},
				},
				Limit:  3,
				Offset: 3,
			},
		},
		{
			scenario: "users_list_partial_filters",
			params: repository.ListUsersParams{
				EmailVerified: new(false),
				BornAfter:     new(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)),
				Sort: []repository.UserSort{
					{Field: repository.SortByFirstName, Desc: true},
					{Field: repository.SortByCreatedAt},
				},
				Limit: 5,
			},
		},
		{
			scenario: "users_list_offset_past_end",
			params: repository.ListUsersParams{
				FirstNamePrefix: new("An"),
				Gender:          new("female"),
				Offset:          50,
			},
		},
	}
	for _, c := range listCases {
		scenario(t, c.scenario, func(t *testing.T, r *run) {
			got, err := record(r, func() (repository.Page[repository.User], error) {
				return r.store.Users().List(t.Context(), c.params)
			})
			if err != nil {
				t.Fatal(err)
			}
			checkResult(t, r, got)
		})
	}

	scenario(t, "users_list_invalid_sort", func(t *testing.T, r *run) {
		_, err := record(r, func() (repository.Page[repository.User], error) {
			return r.store.Users().List(t.Context(), repository.ListUsersParams{
				Sort: []repository.UserSort{{Field: "age"}},
			})
		})
		if !errors.Is(err, repository.ErrInvalidParams) {
			t.Fatalf("got error %v, want %v", err, repository.ErrInvalidParams)
		}
	})

	scenario(t, "users_list_eligible_adults", func(t *testing.T, r *run) {
		got, err := record(r, func() ([]repository.User, error) {
			return r.store.Users().ListEligibleAdults(t.Context())
		})
		if err != nil {
			t.Fatal(err)
		}
		checkResult(t, r, got)
	})

	scenario(t, "users_unverify_minors", func(t *testing.T, r *run) {
		got, err := record(r, func() ([]repository.User, error) {
			return r.store.Users().UnverifyMinors(t.Context())
		})
		if err != nil {
			t.Fatal(err)
		}
		slices.SortFunc(got, func(a, b repository.User) int { return strings.Compare(a.ID, b.ID) })
		checkResult(t, r, got)
	})

	scenario(t, "users_top_age_groups", func(t *testing.T, r *run) {
		got, err := record(r, func() ([]repository.AgeGroup, error) {
			return r.store.Users().TopAgeGroups(t.Context(), 3)
		})
		if err != nil {
			t.Fatal(err)
		}
		checkResult(t, r, got)
	})

	scenario(t, "users_post_stats", func(t *testing.T, r *run) {
		got, err := record(r, func() ([]repository.UserPostStats, error) {
			return r.store.Users().PostStats(t.Context(), 8)
		})
		if err != nil {
			t.Fatal(err)
		}
		checkResult(t, r, got)
	})
}

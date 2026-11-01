package integration

import (
	"testing"

	"github.com/ofadiman/pg-go/seed"
)

func TestSeed(t *testing.T) {
	db := newDB(t)

	t.Run("row counts", func(t *testing.T) {
		tables := []struct {
			name string
			want int
		}{
			{"users", 100},
			{"posts", 405},
			{"comments", 1448},
			{"reactions", 2172},
		}
		for _, table := range tables {
			var got int
			if err := db.QueryRowContext(t.Context(), "select count(*) from "+table.name).Scan(&got); err != nil {
				t.Fatalf("count %s: %v", table.name, err)
			}
			if got != table.want {
				t.Errorf("%s: got %d rows, want %d", table.name, got, table.want)
			}
		}
	})

	t.Run("thread depth", func(t *testing.T) {
		var depth, size int
		err := db.QueryRowContext(t.Context(), `
			with recursive thread as (
				select id, 0 as depth from comments where id = $1
				union all
				select c.id, t.depth + 1 from comments as c join thread as t on c.parent_id = t.id
			)
			select max(depth), count(*) from thread`,
			seed.ThreadRootCommentID,
		).Scan(&depth, &size)
		if err != nil {
			t.Fatalf("query thread: %v", err)
		}
		if depth != 4 || size != 8 {
			t.Errorf("got depth %d with %d comments, want depth 4 with 8 comments", depth, size)
		}
	})

	t.Run("minors", func(t *testing.T) {
		var got int
		err := db.QueryRowContext(t.Context(),
			"select count(*) from users where age($1, birthday) < interval '18 years'",
			now,
		).Scan(&got)
		if err != nil {
			t.Fatalf("count minors: %v", err)
		}
		if got != 8 {
			t.Errorf("got %d minors, want 8", got)
		}
	})

	t.Run("users without activity", func(t *testing.T) {
		var got int
		err := db.QueryRowContext(t.Context(), `
			select count(*)
			from users as u
			where right(u.id::text, 12)::int between 91 and 100
				and (
					exists (select from posts as p where p.user_id = u.id)
					or exists (select from comments as c where c.user_id = u.id)
					or exists (select from reactions as r where r.user_id = u.id)
				)`,
		).Scan(&got)
		if err != nil {
			t.Fatalf("count active users: %v", err)
		}
		if got != 0 {
			t.Errorf("got %d users between 91 and 100 with activity, want 0", got)
		}

		var exists bool
		if err := db.QueryRowContext(t.Context(), "select exists (select from users where id = $1)", seed.UserWithoutActivityID).Scan(&exists); err != nil {
			t.Fatalf("find user without activity: %v", err)
		}
		if !exists {
			t.Errorf("user %s does not exist", seed.UserWithoutActivityID)
		}
	})
}

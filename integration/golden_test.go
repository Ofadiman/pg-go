package integration

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ofadiman/pg-go/repository"
)

var update = flag.Bool("update", false, "rewrite golden files")

type run struct {
	scenario  string
	impl      string
	reference bool
	store     repository.Store
	db        *sql.DB
	rec       *recorder
	queries   []string
	generated []string
}

func scenario(t *testing.T, name string, fn func(t *testing.T, r *run)) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		for i, impl := range implementations {
			t.Run(impl.name, func(t *testing.T) {
				db, rec := newDB(t)
				store, err := impl.new(db, func() time.Time { return now })
				if err != nil {
					t.Fatalf("create store: %v", err)
				}

				r := &run{scenario: name, impl: impl.name, reference: i == 0, store: store, db: db, rec: rec}
				fn(t, r)
				if !t.Failed() {
					checkSQL(t, r)
				}
			})
		}
	})
}

func record[T any](r *run, fn func() (T, error)) (T, error) {
	r.rec.reset()
	value, err := fn()
	r.queries = append(r.queries, r.rec.take()...)
	return value, err
}

func (r *run) record(fn func() error) error {
	_, err := record(r, func() (struct{}, error) { return struct{}{}, fn() })
	return err
}

func (r *run) markGenerated(t *testing.T, ids ...string) {
	t.Helper()
	for _, id := range ids {
		if id == "" {
			t.Fatal("generated id is empty")
		}
		r.generated = append(r.generated, id)
	}
}

func (r *run) replacer() *strings.Replacer {
	pairs := make([]string, 0, len(r.generated)*2)
	for i, id := range r.generated {
		pairs = append(pairs, id, fmt.Sprintf("<generated-%d>", i+1))
	}
	return strings.NewReplacer(pairs...)
}

func checkResult[T any](t *testing.T, r *run, got T) {
	t.Helper()

	data, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}
	data = []byte(r.replacer().Replace(string(data)) + "\n")

	path := filepath.Join("testdata", "results", r.scenario+".json")
	if *update && r.reference {
		writeFile(t, path, data)
		return
	}

	var normalized, want T
	if err := json.Unmarshal(data, &normalized); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if err := json.Unmarshal(readFile(t, path), &want); err != nil {
		t.Fatalf("unmarshal %s: %v", path, err)
	}
	if diff := cmp.Diff(want, normalized, cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("result mismatch (-want +got):\n%s", diff)
	}
}

func checkSQL(t *testing.T, r *run) {
	t.Helper()

	var b strings.Builder
	for i, query := range r.queries {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(strings.TrimSpace(query))
		b.WriteString(";\n")
	}
	got := r.replacer().Replace(b.String())

	path := filepath.Join("testdata", "sql", r.scenario, r.impl+".sql")
	if *update {
		writeFile(t, path, []byte(got))
		return
	}

	if diff := cmp.Diff(string(readFile(t, path)), got); diff != "" {
		t.Errorf("SQL mismatch (-want +got):\n%s", diff)
	}
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden file: %v", err)
	}
	return data
}

func writeFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create golden directory: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write golden file: %v", err)
	}
}

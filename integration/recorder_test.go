package integration

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
)

type recorder struct {
	mu      sync.Mutex
	queries []string
}

func (r *recorder) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.queries = append(r.queries, data.SQL)
	return ctx
}

func (r *recorder) TraceQueryEnd(context.Context, *pgx.Conn, pgx.TraceQueryEndData) {}

func (r *recorder) reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.queries = nil
}

func (r *recorder) take() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	queries := r.queries
	r.queries = nil
	return queries
}

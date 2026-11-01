package goqurepo

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/ofadiman/pg-go/repository"
)

type comments struct {
	*Store
}

func (r comments) Thread(ctx context.Context, rootID string) ([]repository.ThreadComment, error) {
	root := r.q.From("comments").
		Select(goqu.Star(), goqu.L("0").As("depth")).
		Where(goqu.C("id").Eq(rootID))

	replies := r.q.From(goqu.T("comments").As("c")).
		Join(goqu.T("thread").As("t"), goqu.On(goqu.I("c.parent_id").Eq(goqu.I("t.id")))).
		Select(goqu.T("c").All(), goqu.L("? + 1", goqu.I("t.depth")))

	var rows []threadRow
	err := r.q.From("thread").
		WithRecursive("thread", root.UnionAll(replies)).
		Order(goqu.C("depth").Asc(), goqu.C("created_at").Asc(), goqu.C("id").Asc()).
		ScanStructsContext(ctx, &rows)
	if err != nil {
		return nil, fmt.Errorf("get thread %s: %w", rootID, err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("get thread %s: %w", rootID, repository.ErrNotFound)
	}

	out := make([]repository.ThreadComment, len(rows))
	for i, row := range rows {
		out[i] = row.toDomain()
	}
	return out, nil
}

package squirrelrepo

import (
	"context"
	"fmt"

	"github.com/ofadiman/pg-go/repository"
)

type comments struct {
	*Store
}

func (r comments) Thread(ctx context.Context, rootID string) ([]repository.ThreadComment, error) {
	query := r.sb.Select("id", "user_id", "post_id", "parent_id", "created_at", "updated_at", "content", "depth").
		Prefix(
			"WITH RECURSIVE thread AS ("+
				"SELECT *, 0 AS depth FROM comments WHERE id = ? "+
				"UNION ALL "+
				"SELECT c.*, t.depth + 1 FROM comments AS c JOIN thread AS t ON c.parent_id = t.id"+
				")",
			rootID,
		).
		From("thread").
		OrderBy("depth ASC", "created_at ASC", "id ASC")

	var rows []threadRow
	if err := r.selectAll(ctx, &rows, query); err != nil {
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

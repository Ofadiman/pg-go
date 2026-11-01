package seed

import (
	"context"
	"database/sql"
	_ "embed"
)

const (
	FirstUserID           = "00000000-0000-0000-0000-000000000001"
	UserWithoutActivityID = "00000000-0000-0000-0000-000000000100"
	ThreadPostID          = "00000000-0000-0000-0001-000000000001"
	ThreadRootCommentID   = "00000000-0000-0000-0002-900000000001"
)

//go:embed seed.sql
var script string

func Run(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, script)
	return err
}

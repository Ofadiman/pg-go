package goqurepo

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/ofadiman/pg-go/repository"
)

type posts struct {
	*Store
}

func (r posts) CreateMany(ctx context.Context, in []repository.NewPost) ([]repository.Post, error) {
	now := r.now()
	rows := make([]postRow, len(in))
	for i, post := range in {
		rows[i] = newPostRow(post, now)
	}

	var created []postRow
	err := r.q.Insert("posts").
		Rows(rows).
		Returning(goqu.Star()).
		Executor().
		ScanStructsContext(ctx, &created)
	if err != nil {
		return nil, fmt.Errorf("create posts: %w", err)
	}

	out := make([]repository.Post, len(created))
	for i, row := range created {
		out[i] = row.toDomain()
	}
	return out, nil
}

func (r posts) GetWithAuthor(ctx context.Context, id string) (repository.PostWithAuthor, error) {
	var row postWithAuthorRow
	found, err := r.q.From(goqu.T("posts").As("post")).
		Join(goqu.T("users").As("author"), goqu.On(goqu.I("author.id").Eq(goqu.I("post.user_id")))).
		Where(goqu.I("post.id").Eq(id)).
		ScanStructContext(ctx, &row)
	if err != nil {
		return repository.PostWithAuthor{}, fmt.Errorf("get post %s: %w", id, err)
	}
	if !found {
		return repository.PostWithAuthor{}, fmt.Errorf("get post %s: %w", id, repository.ErrNotFound)
	}
	return repository.PostWithAuthor{Post: row.Post.toDomain(), Author: row.Author.toDomain()}, nil
}

func (r posts) ListRecentWithComments(ctx context.Context, limit int) ([]repository.PostWithComments, error) {
	recent := r.q.From("posts").
		Order(goqu.C("created_at").Desc(), goqu.C("id").Asc()).
		Limit(uint(limit))

	var rows []recentPostRow
	err := r.q.From("post").
		With("post", recent).
		Join(goqu.T("users").As("author"), goqu.On(goqu.I("author.id").Eq(goqu.I("post.user_id")))).
		LeftJoin(goqu.T("comments").As("comment"), goqu.On(goqu.I("comment.post_id").Eq(goqu.I("post.id")))).
		LeftJoin(goqu.T("reactions").As("reaction"), goqu.On(goqu.I("reaction.comment_id").Eq(goqu.I("comment.id")))).
		Order(
			goqu.I("post.created_at").Desc(),
			goqu.I("post.id").Asc(),
			goqu.I("comment.created_at").Asc(),
			goqu.I("comment.id").Asc(),
			goqu.I("reaction.created_at").Asc(),
			goqu.I("reaction.id").Asc(),
		).
		ScanStructsContext(ctx, &rows)
	if err != nil {
		return nil, fmt.Errorf("list recent posts: %w", err)
	}

	var out []repository.PostWithComments
	postIndex := make(map[string]int)
	commentIndex := make(map[string]int)
	for _, row := range rows {
		p, ok := postIndex[row.Post.ID]
		if !ok {
			p = len(out)
			postIndex[row.Post.ID] = p
			out = append(out, repository.PostWithComments{Post: row.Post.toDomain(), Author: row.Author.toDomain()})
		}
		if row.Comment.ID == nil {
			continue
		}

		c, ok := commentIndex[*row.Comment.ID]
		if !ok {
			c = len(out[p].Comments)
			commentIndex[*row.Comment.ID] = c
			out[p].Comments = append(out[p].Comments, repository.CommentWithReactions{Comment: row.Comment.toDomain()})
		}
		if row.Reaction.ID == nil {
			continue
		}

		out[p].Comments[c].Reactions = append(out[p].Comments[c].Reactions, row.Reaction.toDomain())
	}
	return out, nil
}

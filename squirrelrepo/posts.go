package squirrelrepo

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/ofadiman/pg-go/repository"
)

type posts struct {
	*Store
}

func (r posts) CreateMany(ctx context.Context, in []repository.NewPost) ([]repository.Post, error) {
	now := r.now()
	query := r.sb.Insert("posts").Columns("user_id", "created_at", "title", "content")
	for _, post := range in {
		query = query.Values(post.UserID, now, post.Title, post.Content)
	}
	query = query.Suffix("RETURNING *")

	var rows []postRow
	if err := r.selectAll(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("create posts: %w", err)
	}

	out := make([]repository.Post, len(rows))
	for i, row := range rows {
		out[i] = row.toDomain()
	}
	return out, nil
}

func (r posts) GetWithAuthor(ctx context.Context, id string) (repository.PostWithAuthor, error) {
	query := r.sb.Select(aliased("post", postColumns)...).
		Columns(aliased("author", userColumns)...).
		From("posts AS post").
		Join("users AS author ON author.id = post.user_id").
		Where(squirrel.Eq{"post.id": id})

	var row postWithAuthorRow
	if err := r.get(ctx, &row, query); err != nil {
		if notFound(err) {
			return repository.PostWithAuthor{}, fmt.Errorf("get post %s: %w", id, repository.ErrNotFound)
		}
		return repository.PostWithAuthor{}, fmt.Errorf("get post %s: %w", id, err)
	}
	return repository.PostWithAuthor{Post: row.Post.toDomain(), Author: row.Author.toDomain()}, nil
}

func (r posts) ListRecentWithComments(ctx context.Context, limit int) ([]repository.PostWithComments, error) {
	recent := squirrel.Select("*").
		From("posts").
		OrderBy("created_at DESC", "id ASC").
		Limit(uint64(limit))

	query := r.sb.Select(aliased("post", postColumns)...).
		Columns(aliased("author", userColumns)...).
		Columns(aliased("comment", commentColumns)...).
		Columns(aliased("reaction", reactionColumns)...).
		PrefixExpr(squirrel.Expr("WITH post AS (?)", recent)).
		From("post").
		Join("users AS author ON author.id = post.user_id").
		LeftJoin("comments AS comment ON comment.post_id = post.id").
		LeftJoin("reactions AS reaction ON reaction.comment_id = comment.id").
		OrderBy(
			"post.created_at DESC",
			"post.id ASC",
			"comment.created_at ASC",
			"comment.id ASC",
			"reaction.created_at ASC",
			"reaction.id ASC",
		)

	var rows []recentPostRow
	if err := r.selectAll(ctx, &rows, query); err != nil {
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

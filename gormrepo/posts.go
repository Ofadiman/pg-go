package gormrepo

import (
	"context"
	"fmt"

	"github.com/ofadiman/pg-go/repository"
	"gorm.io/gorm"
)

type posts struct {
	*Store
}

func (r posts) CreateMany(ctx context.Context, in []repository.NewPost) ([]repository.Post, error) {
	now := r.now()
	rows := make([]postRow, len(in))
	for i, post := range in {
		rows[i] = postRow{UserID: post.UserID, CreatedAt: now, Title: post.Title, Content: post.Content}
	}

	if err := r.q.WithContext(ctx).Create(&rows).Error; err != nil {
		return nil, fmt.Errorf("create posts: %w", err)
	}

	out := make([]repository.Post, len(rows))
	for i, row := range rows {
		out[i] = row.toDomain()
	}
	return out, nil
}

func (r posts) GetWithAuthor(ctx context.Context, id string) (repository.PostWithAuthor, error) {
	var row postRow
	if err := r.q.WithContext(ctx).Joins("Author").First(&row, "posts.id = ?", id).Error; err != nil {
		if notFound(err) {
			return repository.PostWithAuthor{}, fmt.Errorf("get post %s: %w", id, repository.ErrNotFound)
		}
		return repository.PostWithAuthor{}, fmt.Errorf("get post %s: %w", id, err)
	}
	return repository.PostWithAuthor{Post: row.toDomain(), Author: row.Author.toDomain()}, nil
}

func (r posts) ListRecentWithComments(ctx context.Context, limit int) ([]repository.PostWithComments, error) {
	chronological := func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC, id ASC")
	}

	var rows []postRow
	err := r.q.WithContext(ctx).
		Joins("Author").
		Preload("Comments", chronological).
		Preload("Comments.Reactions", chronological).
		Order("posts.created_at DESC, posts.id ASC").
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list recent posts: %w", err)
	}

	out := make([]repository.PostWithComments, len(rows))
	for i, row := range rows {
		post := repository.PostWithComments{Post: row.toDomain(), Author: row.Author.toDomain()}
		for _, comment := range row.Comments {
			withReactions := repository.CommentWithReactions{Comment: comment.toDomain()}
			for _, reaction := range comment.Reactions {
				withReactions.Reactions = append(withReactions.Reactions, reaction.toDomain())
			}
			post.Comments = append(post.Comments, withReactions)
		}
		out[i] = post
	}
	return out, nil
}

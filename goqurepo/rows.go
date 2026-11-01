package goqurepo

import (
	"time"

	"github.com/ofadiman/pg-go/repository"
)

type userRow struct {
	ID            string             `db:"id" goqu:"skipinsert"`
	CreatedAt     time.Time          `db:"created_at"`
	UpdatedAt     *time.Time         `db:"updated_at"`
	Address       repository.Address `db:"address"`
	Birthday      time.Time          `db:"birthday"`
	Email         string             `db:"email"`
	EmailVerified bool               `db:"email_verified"`
	FirstName     string             `db:"first_name"`
	Gender        string             `db:"gender"`
	LastName      string             `db:"last_name"`
}

func newUserRow(in repository.NewUser, now time.Time) userRow {
	return userRow{
		CreatedAt:     now,
		Address:       in.Address,
		Birthday:      in.Birthday,
		Email:         in.Email,
		EmailVerified: in.EmailVerified,
		FirstName:     in.FirstName,
		Gender:        in.Gender,
		LastName:      in.LastName,
	}
}

func (r userRow) toDomain() repository.User {
	return repository.User{
		ID:            r.ID,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
		Address:       r.Address,
		Birthday:      r.Birthday,
		Email:         r.Email,
		EmailVerified: r.EmailVerified,
		FirstName:     r.FirstName,
		Gender:        r.Gender,
		LastName:      r.LastName,
	}
}

func usersToDomain(rows []userRow) []repository.User {
	out := make([]repository.User, len(rows))
	for i, row := range rows {
		out[i] = row.toDomain()
	}
	return out
}

type postRow struct {
	ID        string     `db:"id" goqu:"skipinsert"`
	UserID    string     `db:"user_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	Title     string     `db:"title"`
	Content   string     `db:"content"`
}

func newPostRow(in repository.NewPost, now time.Time) postRow {
	return postRow{
		UserID:    in.UserID,
		CreatedAt: now,
		Title:     in.Title,
		Content:   in.Content,
	}
}

func (r postRow) toDomain() repository.Post {
	return repository.Post{
		ID:        r.ID,
		UserID:    r.UserID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		Title:     r.Title,
		Content:   r.Content,
	}
}

type postWithAuthorRow struct {
	Post   postRow `db:"post"`
	Author userRow `db:"author"`
}

type recentPostRow struct {
	Post     postRow             `db:"post"`
	Author   userRow             `db:"author"`
	Comment  nullableCommentRow  `db:"comment"`
	Reaction nullableReactionRow `db:"reaction"`
}

type nullableCommentRow struct {
	ID        *string    `db:"id"`
	UserID    *string    `db:"user_id"`
	PostID    *string    `db:"post_id"`
	ParentID  *string    `db:"parent_id"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	Content   *string    `db:"content"`
}

func (r nullableCommentRow) toDomain() repository.Comment {
	return repository.Comment{
		ID:        *r.ID,
		UserID:    *r.UserID,
		PostID:    *r.PostID,
		ParentID:  r.ParentID,
		CreatedAt: *r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		Content:   *r.Content,
	}
}

type nullableReactionRow struct {
	ID        *string    `db:"id"`
	UserID    *string    `db:"user_id"`
	CommentID *string    `db:"comment_id"`
	CreatedAt *time.Time `db:"created_at"`
	Emoji     *string    `db:"emoji"`
}

func (r nullableReactionRow) toDomain() repository.Reaction {
	return repository.Reaction{
		ID:        *r.ID,
		UserID:    *r.UserID,
		CommentID: *r.CommentID,
		CreatedAt: *r.CreatedAt,
		Emoji:     *r.Emoji,
	}
}

type threadRow struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	PostID    string     `db:"post_id"`
	ParentID  *string    `db:"parent_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	Content   string     `db:"content"`
	Depth     int        `db:"depth"`
}

func (r threadRow) toDomain() repository.ThreadComment {
	return repository.ThreadComment{
		Comment: repository.Comment{
			ID:        r.ID,
			UserID:    r.UserID,
			PostID:    r.PostID,
			ParentID:  r.ParentID,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
			Content:   r.Content,
		},
		Depth: r.Depth,
	}
}

type ageGroupRow struct {
	Age   int `db:"age"`
	Total int `db:"total"`
}

type userPostStatsRow struct {
	UserID       string    `db:"user_id"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	PostCount    int       `db:"post_count"`
	CommentCount int       `db:"comment_count"`
	LastPostAt   time.Time `db:"last_post_at"`
}

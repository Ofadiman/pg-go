package gormrepo

import (
	"time"

	"github.com/ofadiman/pg-go/repository"
)

type userRow struct {
	ID            string `gorm:"primaryKey;default:uuid_generate_v4()"`
	CreatedAt     time.Time
	UpdatedAt     *time.Time `gorm:"autoUpdateTime:false"`
	Address       repository.Address
	Birthday      time.Time
	Email         string
	EmailVerified bool
	FirstName     string
	Gender        string
	LastName      string
}

func (userRow) TableName() string {
	return "users"
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
	ID        string `gorm:"primaryKey;default:uuid_generate_v4()"`
	UserID    string
	CreatedAt time.Time
	UpdatedAt *time.Time `gorm:"autoUpdateTime:false"`
	Title     string
	Content   string
	Author    *userRow     `gorm:"foreignKey:UserID"`
	Comments  []commentRow `gorm:"foreignKey:PostID"`
}

func (postRow) TableName() string {
	return "posts"
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

type commentRow struct {
	ID        string `gorm:"primaryKey"`
	UserID    string
	PostID    string
	ParentID  *string
	CreatedAt time.Time
	UpdatedAt *time.Time `gorm:"autoUpdateTime:false"`
	Content   string
	Reactions []reactionRow `gorm:"foreignKey:CommentID"`
}

func (commentRow) TableName() string {
	return "comments"
}

func (r commentRow) toDomain() repository.Comment {
	return repository.Comment{
		ID:        r.ID,
		UserID:    r.UserID,
		PostID:    r.PostID,
		ParentID:  r.ParentID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		Content:   r.Content,
	}
}

type reactionRow struct {
	ID        string `gorm:"primaryKey"`
	UserID    string
	CommentID string
	CreatedAt time.Time
	Emoji     string
}

func (reactionRow) TableName() string {
	return "reactions"
}

func (r reactionRow) toDomain() repository.Reaction {
	return repository.Reaction{
		ID:        r.ID,
		UserID:    r.UserID,
		CommentID: r.CommentID,
		CreatedAt: r.CreatedAt,
		Emoji:     r.Emoji,
	}
}

type threadRow struct {
	ID        string
	UserID    string
	PostID    string
	ParentID  *string
	CreatedAt time.Time
	UpdatedAt *time.Time
	Content   string
	Depth     int
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
	Age   int
	Total int
}

type userPostStatsRow struct {
	UserID       string
	FirstName    string
	LastName     string
	PostCount    int
	CommentCount int
	LastPostAt   time.Time
}

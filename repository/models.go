package repository

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Address struct {
	City         string `json:"city"`
	Country      string `json:"country"`
	Street       string `json:"street"`
	StreetNumber string `json:"street_number"`
	ZipCode      string `json:"zip_code"`
}

func (a Address) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Address) Scan(value any) error {
	switch v := value.(type) {
	case nil:
		return nil
	case []byte:
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	default:
		return fmt.Errorf("scan address: unsupported type %T", value)
	}
}

type User struct {
	ID            string
	CreatedAt     time.Time
	UpdatedAt     *time.Time
	Address       Address
	Birthday      time.Time
	Email         string
	EmailVerified bool
	FirstName     string
	Gender        string
	LastName      string
}

type NewUser struct {
	Address       Address
	Birthday      time.Time
	Email         string
	EmailVerified bool
	FirstName     string
	Gender        string
	LastName      string
}

type UserUpdate struct {
	FirstName     *string
	LastName      *string
	EmailVerified *bool
}

type Post struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	UpdatedAt *time.Time
	Title     string
	Content   string
}

type NewPost struct {
	UserID  string
	Title   string
	Content string
}

type Comment struct {
	ID        string
	UserID    string
	PostID    string
	ParentID  *string
	CreatedAt time.Time
	UpdatedAt *time.Time
	Content   string
}

type Reaction struct {
	ID        string
	UserID    string
	CommentID string
	CreatedAt time.Time
	Emoji     string
}

type PostWithAuthor struct {
	Post   Post
	Author User
}

type PostWithComments struct {
	Post     Post
	Author   User
	Comments []CommentWithReactions
}

type CommentWithReactions struct {
	Comment   Comment
	Reactions []Reaction
}

type ThreadComment struct {
	Comment Comment
	Depth   int
}

type AgeGroup struct {
	Age   int
	Total int
}

type UserPostStats struct {
	UserID       string
	FirstName    string
	LastName     string
	PostCount    int
	CommentCount int
	LastPostAt   time.Time
}

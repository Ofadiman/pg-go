package goqu

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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
	if value == nil {
		return nil
	}

	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return errors.New("cannot scan into Address: unsupported type")
	}

	return json.Unmarshal(b, a)
}

type User struct {
	ID string `db:"id" goqu:"skipinsert"`

	CreatedAt time.Time  `db:"created_at" goqu:"skipinsert"`
	UpdatedAt *time.Time `db:"updated_at"`

	Address       Address   `db:"address"`
	Birthday      time.Time `db:"birthday"`
	Email         string    `db:"email"`
	EmailVerified bool      `db:"email_verified"`
	FirstName     string    `db:"first_name"`
	Gender        string    `db:"gender"`
	LastName      string    `db:"last_name"`
}

type Post struct {
	ID     string `db:"id" goqu:"skipinsert"`
	UserID string `db:"user_id"`

	CreatedAt time.Time  `db:"created_at" goqu:"skipinsert"`
	UpdatedAt *time.Time `db:"updated_at"`

	Title   string `db:"title"`
	Content string `db:"content"`
}

type UserIDs []string

func (u *UserIDs) Scan(value any) error {
	if value == nil {
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}

	return json.Unmarshal(data, u)
}

func (u UserIDs) Value() (driver.Value, error) {
	return json.Marshal(u)
}

type UserIdsByCountry struct {
	Country string  `db:"country"`
	Total   int     `db:"total"`
	UserIDs UserIDs `db:"user_ids"`
}

type MostPopularAgeGroup struct {
	Total int `db:"total"`
	Age   int `db:"age"`
}

type RecentPostWithRelatedRecords struct {
	PostID        string     `db:"post_id"`
	PostUserID    string     `db:"post_user_id"`
	PostCreatedAt time.Time  `db:"post_created_at"`
	PostUpdatedAt *time.Time `db:"post_updated_at"`
	PostContent   string     `db:"post_content"`
	PostTitle     string     `db:"post_title"`

	UserID            string     `db:"user_id"`
	UserCreatedAt     time.Time  `db:"user_created_at"`
	UserUpdatedAt     *time.Time `db:"user_updated_at"`
	UserAddress       Address    `db:"user_address"`
	UserBirthday      time.Time  `db:"user_birthday"`
	UserEmail         string     `db:"user_email"`
	UserEmailVerified bool       `db:"user_email_verified"`
	UserFirstName     string     `db:"user_first_name"`
	UserGender        string     `db:"user_gender"`
	UserLastName      string     `db:"user_last_name"`

	CommentID        *string    `db:"comment_id"`
	CommentUserID    *string    `db:"comment_user_id"`
	CommentParentID  *string    `db:"comment_parent_id"`
	CommentPostID    *string    `db:"comment_post_id"`
	CommentCreatedAt *time.Time `db:"comment_created_at"`
	CommentUpdatedAt *time.Time `db:"comment_updated_at"`
	CommentContent   *string    `db:"comment_content"`

	ReactionID        *string    `db:"reaction_id"`
	ReactionCommentID *string    `db:"reaction_comment_id"`
	ReactionUserID    *string    `db:"reaction_user_id"`
	ReactionCreatedAt *time.Time `db:"reaction_created_at"`
	ReactionEmoji     *string    `db:"reaction_emoji"`
}

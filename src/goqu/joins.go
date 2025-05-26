package main

import (
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ofadiman/pg-go/src/utils"
)

type PostData struct {
	ID        string
	AuthorID  string
	CreatedAt time.Time
	UpdatedAt *time.Time
	Title     string
	Content   string
}

type UserData struct {
	ID            string
	CreatedAt     time.Time
	UpdatedAt     *time.Time
	FirstName     string
	LastName      string
	Email         string
	EmailVerified bool
	PhoneNumber   string
	Gender        string
	City          string
	Country       string
	Street        string
	StreetNumber  string
	ZipCode       string
	Birthday      time.Time
	Posts         []PostData
}

type UserPostRow struct {
	UsersID            string     `db:"users_id"`
	UsersCreatedAt     time.Time  `db:"users_created_at"`
	UsersUpdatedAt     *time.Time `db:"users_updated_at"`
	UsersFirstName     string     `db:"users_first_name"`
	UsersLastName      string     `db:"users_last_name"`
	UsersEmail         string     `db:"users_email"`
	UsersEmailVerified bool       `db:"users_email_verified"`
	UsersPhoneNumber   string     `db:"users_phone_number"`
	UsersGender        string     `db:"users_gender"`
	UsersCity          string     `db:"users_city"`
	UsersCountry       string     `db:"users_country"`
	UsersStreet        string     `db:"users_street"`
	UsersStreetNumber  string     `db:"users_street_number"`
	UsersZipCode       string     `db:"users_zip_code"`
	UsersBirthday      time.Time  `db:"users_birthday"`

	PostsID        *string    `db:"posts_id"`
	PostsAuthorID  *string    `db:"posts_author_id"`
	PostsCreatedAt *time.Time `db:"posts_created_at"`
	PostsUpdatedAt *time.Time `db:"posts_updated_at"`
	PostsTitle     *string    `db:"posts_title"`
	PostsContent   *string    `db:"posts_content"`
}

func Joins() {
	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	database := goqu.New("postgres", sqlxdb.DB)

	statement := database.From("users").
		LeftJoin(goqu.T("posts"), goqu.On(goqu.Ex{"posts.author_id": goqu.I("users.id")})).
		Select(
			goqu.I("users.id").As("users_id"),
			goqu.I("users.created_at").As("users_created_at"),
			goqu.I("users.updated_at").As("users_updated_at"),
			goqu.I("users.first_name").As("users_first_name"),
			goqu.I("users.last_name").As("users_last_name"),
			goqu.I("users.email").As("users_email"),
			goqu.I("users.email_verified").As("users_email_verified"),
			goqu.I("users.phone_number").As("users_phone_number"),
			goqu.I("users.gender").As("users_gender"),
			goqu.I("users.city").As("users_city"),
			goqu.I("users.country").As("users_country"),
			goqu.I("users.street").As("users_street"),
			goqu.I("users.street_number").As("users_street_number"),
			goqu.I("users.zip_code").As("users_zip_code"),
			goqu.I("users.birthday").As("users_birthday"),

			goqu.I("posts.id").As("posts_id"),
			goqu.I("posts.author_id").As("posts_author_id"),
			goqu.I("posts.created_at").As("posts_created_at"),
			goqu.I("posts.updated_at").As("posts_updated_at"),
			goqu.I("posts.title").As("posts_title"),
			goqu.I("posts.content").As("posts_content"),
		)

	var rows []UserPostRow
	err = statement.Executor().ScanStructs(&rows)
	if err != nil {
		panic(err)
	}

	usersMap := make(map[string]*UserData)
	for _, row := range rows {
		user, exists := usersMap[row.UsersID]
		if !exists {
			user = &UserData{
				ID:            row.UsersID,
				CreatedAt:     row.UsersCreatedAt,
				UpdatedAt:     row.UsersUpdatedAt,
				FirstName:     row.UsersFirstName,
				LastName:      row.UsersLastName,
				Email:         row.UsersEmail,
				EmailVerified: row.UsersEmailVerified,
				PhoneNumber:   row.UsersPhoneNumber,
				Gender:        row.UsersGender,
				City:          row.UsersCity,
				Country:       row.UsersCountry,
				Street:        row.UsersStreet,
				StreetNumber:  row.UsersStreetNumber,
				ZipCode:       row.UsersZipCode,
				Birthday:      row.UsersBirthday,
				Posts:         []PostData{},
			}
			usersMap[row.UsersID] = user
		}

		if row.PostsID != nil {
			user.Posts = append(user.Posts, PostData{
				ID:        *row.PostsID,
				AuthorID:  *row.PostsAuthorID,
				CreatedAt: *row.PostsCreatedAt,
				UpdatedAt: row.PostsUpdatedAt,
				Title:     *row.PostsTitle,
				Content:   *row.PostsContent,
			})
		}
	}

	var users []UserData
	for _, user := range usersMap {
		users = append(users, *user)
	}

	utils.Print(users)
}

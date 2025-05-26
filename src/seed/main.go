package main

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ofadiman/pg-go/src/utils"
)

func init() {
	gofakeit.Seed(1)
}

type User struct {
	ID string `db:"id"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`

	FirstName     string    `db:"first_name"`
	LastName      string    `db:"last_name"`
	Email         string    `db:"email"`
	EmailVerified bool      `db:"email_verified"`
	PhoneNumber   string    `db:"phone_number"`
	Gender        string    `db:"gender"`
	City          string    `db:"city"`
	Country       string    `db:"country"`
	Street        string    `db:"street"`
	StreetNumber  string    `db:"street_number"`
	ZipCode       string    `db:"zip_code"`
	Birthday      time.Time `db:"birthday"`
}

type Post struct {
	ID       string `db:"id"`
	AuthorID string `db:"author_id"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`

	Title   string `db:"title"`
	Content string `db:"content"`
}

func main() {
	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	sqlxdb.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")

	database := goqu.New("postgres", sqlxdb.DB)

	var users []User

	for range 100 {
		firstName := gofakeit.FirstName()
		lastName := gofakeit.LastName()
		email := fmt.Sprintf("%s.%s@%s.%s", firstName, lastName, gofakeit.DomainName(), gofakeit.DomainSuffix())
		createdAt := gofakeit.DateRange(
			time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 30, 0, 0, 0, 0, time.UTC),
		)
		updatedAt := createdAt.Add(time.Hour * 24 * time.Duration(gofakeit.IntN(365)))

		user := User{
			ID:            gofakeit.UUID(),
			CreatedAt:     createdAt,
			UpdatedAt:     &updatedAt,
			FirstName:     firstName,
			LastName:      lastName,
			Email:         email,
			EmailVerified: gofakeit.Bool(),
			PhoneNumber:   gofakeit.PhoneFormatted(),
			Gender:        gofakeit.Gender(),
			City:          gofakeit.City(),
			Country:       gofakeit.Country(),
			Street:        gofakeit.Street(),
			StreetNumber:  gofakeit.StreetNumber(),
			ZipCode:       gofakeit.Zip(),
			Birthday: gofakeit.DateRange(
				time.Date(1950, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC),
			),
		}

		users = append(users, user)
	}

	database.Insert("users").Rows(users).Executor().Exec()

	var posts []Post

	for _, user := range users {
		postsCount := gofakeit.IntN(10)
		for range postsCount {
			createdAt := user.CreatedAt.Add(time.Hour * 24 * time.Duration(gofakeit.IntN(365)))
			var updatedAt *time.Time
			if gofakeit.Bool() {
				updatedAt = utils.Pointer(createdAt.Add(time.Hour * 24 * time.Duration(gofakeit.IntN(30))))
			}

			post := Post{
				ID:        gofakeit.UUID(),
				AuthorID:  user.ID,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				Title:     gofakeit.LoremIpsumSentence(gofakeit.IntN(9) + 1),
				Content:   gofakeit.LoremIpsumParagraph(1, gofakeit.IntN(4)+1, gofakeit.IntN(9)+1, " "),
			}

			posts = append(posts, post)
		}
	}

	database.Insert("posts").Rows(posts).Executor().Exec()
}

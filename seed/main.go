package seed

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ofadiman/pg-go/src/utils"
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
	ID string `db:"id"`

	CreatedAt time.Time  `db:"created_at"`
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
	ID       string `db:"id"`
	AuthorID string `db:"author_id"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`

	Title   string `db:"title"`
	Content string `db:"content"`
}

type Comment struct {
	ID       string  `db:"id"`
	AuthorID string  `db:"author_id"`
	PostID   string  `db:"post_id"`
	ParentID *string `db:"parent_id"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`

	Content string `db:"content"`
}

type Reaction struct {
	ID        string `db:"id"`
	AuthorID  string `db:"author_id"`
	CommentID string `db:"comment_id"`

	CreatedAt time.Time `db:"created_at"`

	Emoji string `db:"emoji"`
}

var emojis = []string{
	"face_with_tears_of_joy",
	"red_heart",
	"smiling_face_with_heart_eyes",
	"thumbs_up",
	"smiling_face_with_smiling_eyes",
	"fire",
	"loudly_crying_face",
	"hundred_points",
	"party_popper",
	"clapping_hands",
}

func Seed() {
	gofakeit.Seed(1)

	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	tx, err := sqlxdb.Beginx()
	if err != nil {
		panic(err)
	}

	_, err = tx.Exec("TRUNCATE TABLE users CASCADE")
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()

	db := goqu.New("postgres", sqlxdb.DB)

	var users []User

	for range 100 {
		firstName := gofakeit.FirstName()
		lastName := gofakeit.LastName()
		email := fmt.Sprintf("%[1]s.%[2]s@%[3]s.%[4]s", firstName, lastName, gofakeit.DomainName(), gofakeit.DomainSuffix())
		createdAt := gofakeit.DateRange(
			time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 30, 0, 0, 0, 0, time.UTC),
		)
		updatedAt := createdAt.Add(time.Hour * 24 * time.Duration(gofakeit.IntN(365)))

		user := User{
			ID:        gofakeit.UUID(),
			CreatedAt: createdAt,
			UpdatedAt: &updatedAt,
			Address: Address{
				City:         gofakeit.City(),
				Country:      gofakeit.Country(),
				Street:       gofakeit.Street(),
				StreetNumber: gofakeit.StreetNumber(),
				ZipCode:      gofakeit.Zip(),
			},
			FirstName:     firstName,
			LastName:      lastName,
			Email:         email,
			EmailVerified: gofakeit.Bool(),
			Gender:        gofakeit.Gender(),
			Birthday: gofakeit.DateRange(
				time.Date(1950, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC),
			),
		}

		users = append(users, user)
	}

	_, err = db.Insert("users").Rows(users).Executor().Exec()
	if err != nil {
		panic(err)
	}

	var posts []Post

	for _, user := range users {
		postsCount := gofakeit.IntN(10)

		for range postsCount {
			post := Post{
				ID:       gofakeit.UUID(),
				AuthorID: user.ID,
				Title:    gofakeit.LoremIpsumSentence(gofakeit.IntN(9) + 1),
				Content:  gofakeit.LoremIpsumParagraph(1, gofakeit.IntN(4)+1, gofakeit.IntN(9)+1, " "),
			}

			post.CreatedAt = user.CreatedAt.Add(time.Hour * 24 * time.Duration(gofakeit.IntN(365)))

			if gofakeit.Bool() {
				post.UpdatedAt = utils.Pointer(post.CreatedAt.Add(time.Hour * 24 * time.Duration(gofakeit.IntN(30))))
			}

			posts = append(posts, post)
		}
	}

	_, err = db.Insert("posts").Rows(posts).Executor().Exec()
	if err != nil {
		panic(err)
	}

	var comments []Comment

	for _, post := range posts {
		commentsCount := gofakeit.IntN(10)

		for range commentsCount {
			comment := Comment{
				ID:       gofakeit.UUID(),
				PostID:   post.ID,
				AuthorID: users[gofakeit.IntN(len(users)-1)].ID,
				Content:  gofakeit.LoremIpsumParagraph(1, gofakeit.IntN(4)+1, gofakeit.IntN(9)+1, " "),
			}

			comment.CreatedAt = post.CreatedAt.Add(time.Hour * 24 * time.Duration(gofakeit.IntN(30)))

			if gofakeit.Bool() {
				comment.UpdatedAt = utils.Pointer(comment.CreatedAt.Add(time.Hour * 24 * time.Duration(gofakeit.IntN(30))))
			}

			comments = append(comments, comment)

			if gofakeit.Bool() {
				nestedCommentsCount := gofakeit.IntN(5)

				for range nestedCommentsCount {
					nestedComment := Comment{
						ID:       gofakeit.UUID(),
						PostID:   post.ID,
						AuthorID: users[gofakeit.IntN(len(users)-1)].ID,
						ParentID: &comment.ID,
						Content:  gofakeit.LoremIpsumParagraph(1, gofakeit.IntN(4)+1, gofakeit.IntN(9)+1, " "),
					}

					nestedComment.CreatedAt = comment.CreatedAt.Add(time.Hour * 24 * time.Duration(gofakeit.IntN(5)))

					if gofakeit.Bool() {
						nestedComment.UpdatedAt = utils.Pointer(nestedComment.CreatedAt.Add(time.Hour * 24 * time.Duration(gofakeit.IntN(3))))
					}

					comments = append(comments, nestedComment)
				}
			}
		}
	}

	_, err = db.Insert("comments").Rows(comments).Executor().Exec()
	if err != nil {
		panic(err)
	}

	var reactions []Reaction

	for _, comment := range comments {
		reactionsCount := gofakeit.IntN(10)

		for range reactionsCount {
			reaction := Reaction{
				ID:        gofakeit.UUID(),
				AuthorID:  users[gofakeit.IntN(len(users)-1)].ID,
				CommentID: comment.ID,
				CreatedAt: comment.CreatedAt.Add(time.Hour * 24 * time.Duration(gofakeit.IntN(10))),
				Emoji:     emojis[gofakeit.IntN(len(emojis)-1)],
			}

			reactions = append(reactions, reaction)
		}
	}

	_, err = db.Insert("reactions").Rows(reactions).Executor().Exec()
	if err != nil {
		panic(err)
	}
}

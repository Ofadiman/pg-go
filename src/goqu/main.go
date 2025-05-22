package main

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ofadiman/pg-go/src/utils"
)

type User struct {
	ID string `db:"id" goqu:"skipinsert"`

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

type UserRepository struct {
	db *goqu.Database
}

func createUser() User {
	return User{
		FirstName:     gofakeit.FirstName(),
		LastName:      gofakeit.LastName(),
		Email:         gofakeit.Email(),
		EmailVerified: gofakeit.Bool(),
		PhoneNumber:   gofakeit.PhoneFormatted(),
		Gender:        gofakeit.Gender(),
		City:          gofakeit.City(),
		Country:       gofakeit.Country(),
		Street:        gofakeit.Street(),
		StreetNumber:  gofakeit.StreetNumber(),
		ZipCode:       gofakeit.Zip(),
		Birthday:      gofakeit.PastDate(),
		CreatedAt:     gofakeit.PastDate(),
	}
}

func (r *UserRepository) insertOne() User {
	user := createUser()
	var inserted User

	ds := r.db.Insert("users").Rows(user).Returning(goqu.Star())
	_, err := ds.Executor().ScanStruct(&inserted)
	if err != nil {
		panic(err)
	}
	return inserted
}

func (r *UserRepository) insertMany() []User {
	users := []User{createUser(), createUser(), createUser()}
	var inserted []User

	ds := r.db.Insert("users").Rows(users).Returning(goqu.Star())
	err := ds.Executor().ScanStructs(&inserted)
	if err != nil {
		panic(err)
	}
	return inserted
}

func main() {
	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	db := goqu.New("postgres", sqlxdb.DB)

	repository := UserRepository{db: db}

	user := repository.insertOne()
	utils.Print(user)

	users := repository.insertMany()
	utils.Print(users)
}

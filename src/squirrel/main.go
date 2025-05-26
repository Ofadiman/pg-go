package main

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ofadiman/pg-go/src/utils"
)

func init() {
	gofakeit.Seed(375)
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

type UserRepository struct {
	db *sqlx.DB
}

func (this *UserRepository) insertOne() User {
	sql, args, err := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Insert("users").
		Columns("id", "first_name", "last_name", "email", "email_verified", "phone_number", "gender", "city", "country", "street", "street_number", "zip_code", "birthday", "created_at").
		Values(gofakeit.UUID(), gofakeit.FirstName(), gofakeit.LastName(), gofakeit.Email(), gofakeit.Bool(), gofakeit.PhoneFormatted(), gofakeit.Gender(), gofakeit.City(), gofakeit.Country(), gofakeit.Street(), gofakeit.StreetNumber(), gofakeit.Zip(), gofakeit.PastDate(), gofakeit.PastDate()).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		panic(err)
	}

	var user User
	err = this.db.Get(&user, sql, args...)
	if err != nil {
		panic(err)
	}

	return user
}

func (this *UserRepository) insertMany() []User {
	sql, args, err := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Insert("users").
		Columns("id", "first_name", "last_name", "email", "email_verified", "phone_number", "gender", "city", "country", "street", "street_number", "zip_code", "birthday", "created_at").
		Values(gofakeit.UUID(), gofakeit.FirstName(), gofakeit.LastName(), gofakeit.Email(), gofakeit.Bool(), gofakeit.PhoneFormatted(), gofakeit.Gender(), gofakeit.City(), gofakeit.Country(), gofakeit.Street(), gofakeit.StreetNumber(), gofakeit.Zip(), gofakeit.PastDate(), gofakeit.PastDate()).
		Values(gofakeit.UUID(), gofakeit.FirstName(), gofakeit.LastName(), gofakeit.Email(), gofakeit.Bool(), gofakeit.PhoneFormatted(), gofakeit.Gender(), gofakeit.City(), gofakeit.Country(), gofakeit.Street(), gofakeit.StreetNumber(), gofakeit.Zip(), gofakeit.PastDate(), gofakeit.PastDate()).
		Values(gofakeit.UUID(), gofakeit.FirstName(), gofakeit.LastName(), gofakeit.Email(), gofakeit.Bool(), gofakeit.PhoneFormatted(), gofakeit.Gender(), gofakeit.City(), gofakeit.Country(), gofakeit.Street(), gofakeit.StreetNumber(), gofakeit.Zip(), gofakeit.PastDate(), gofakeit.PastDate()).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		panic(err)
	}

	var users []User
	err = this.db.Select(&users, sql, args...)
	if err != nil {
		panic(err)
	}

	return users
}

func (this *UserRepository) selectById() User {
	sql, args, err := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("*").
		From("users").
		Where(squirrel.Eq{"id": "8917f080-70ea-4c89-b86c-859d617f3a0b"}).
		ToSql()
	if err != nil {
		panic(err)
	}

	var user User
	err = this.db.Get(&user, sql, args...)
	if err != nil {
		panic(err)
	}

	return user
}

func (this *UserRepository) paginate() []User {
	statement := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("*").
		From("users")
	statement = statement.Where(squirrel.Or{
		squirrel.Eq{"id": "49b7cf84-d525-4456-af6e-5b4015473cf6"},
		squirrel.Eq{"first_name": []string{"Kirsten", "Coleman"}},
	})
	statement = statement.Limit(100)
	statement = statement.Offset(0)

	sql, args, err := statement.ToSql()
	if err != nil {
		panic(err)
	}

	var user []User
	err = this.db.Select(&user, sql, args...)
	if err != nil {
		panic(err)
	}

	return user
}

func main() {
	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	sqlxdb.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")

	repository := UserRepository{db: sqlxdb}

	user := repository.insertOne()
	utils.Print(user)

	users := repository.insertMany()
	utils.Print(users)

	userById := repository.selectById()
	utils.Print(userById)

	paginated := repository.paginate()
	utils.Print(paginated)
}

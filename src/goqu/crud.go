package main

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ofadiman/pg-go/src/utils"
)

func init() {
	gofakeit.Seed(1)
}

type UserRepository struct {
	db *goqu.Database
}

func createUser() User {
	return User{
		ID:            gofakeit.UUID(),
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

func (this *UserRepository) insertOne() User {
	user := createUser()
	var inserted User

	statement := this.db.Insert("users").Rows(user).Returning(goqu.Star())
	_, err := statement.Executor().ScanStruct(&inserted)
	if err != nil {
		panic(err)
	}
	return inserted
}

func (this *UserRepository) insertMany() []User {
	users := []User{createUser(), createUser(), createUser()}
	var inserted []User

	statement := this.db.Insert("users").Rows(users).Returning(goqu.Star())
	err := statement.Executor().ScanStructs(&inserted)
	if err != nil {
		panic(err)
	}
	return inserted
}

func (this *UserRepository) selectById() User {
	var user User
	statement := this.db.From("users").
		Where(goqu.C("id").Eq("b80bacdc-c556-426d-b31c-2dea8f7eed48"))
	_, err := statement.Executor().ScanStruct(&user)
	if err != nil {
		panic(err)
	}
	return user
}

func (this *UserRepository) paginate() []User {
	var users []User
	statement := this.db.From("users")
	statement = statement.Where(goqu.Or(
		goqu.C("id").Eq("064be158-eb2c-41dc-9308-8585917cafee"),
		goqu.C("first_name").In([]string{"Rubye", "Myron"}),
	))
	statement = statement.Limit(100)
	statement = statement.Offset(0)

	err := statement.Executor().ScanStructs(&users)
	if err != nil {
		panic(err)
	}
	return users
}

func (this *UserRepository) updateById() User {
	var user User
	statement := this.db.Update("users").
		Set(goqu.Record{
			"first_name": "Jack",
			"last_name":  "Sparrow",
			"updated_at": goqu.Literal("NOW()"),
		}).
		Where(goqu.C("id").Eq("b80bacdc-c556-426d-b31c-2dea8f7eed48")).
		Returning(goqu.Star())

	_, err := statement.Executor().ScanStruct(&user)
	if err != nil {
		panic(err)
	}
	return user
}

func (this *UserRepository) updateAll() []User {
	var users []User
	statement := this.db.Update("users").
		Set(goqu.Record{
			"updated_at": goqu.Literal("NOW()"),
		}).
		Returning(goqu.Star())

	err := statement.Executor().ScanStructs(&users)
	if err != nil {
		panic(err)
	}
	return users
}

func (this *UserRepository) deleteById() User {
	var user User
	statement := this.db.Delete("users").
		Where(goqu.C("id").Eq("064be158-eb2c-41dc-9308-8585917cafee")).
		Returning(goqu.Star())

	_, err := statement.Executor().ScanStruct(&user)
	if err != nil {
		panic(err)
	}
	return user
}

func (this *UserRepository) deleteAll() []User {
	var users []User
	statement := this.db.Delete("users").
		Returning(goqu.Star())

	err := statement.Executor().ScanStructs(&users)
	if err != nil {
		panic(err)
	}
	return users
}

func Crud() {
	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	sqlxdb.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")

	db := goqu.New("postgres", sqlxdb.DB)

	repository := UserRepository{db: db}

	user := repository.insertOne()
	utils.Print(user)

	users := repository.insertMany()
	utils.Print(users)

	userById := repository.selectById()
	utils.Print(userById)

	paginated := repository.paginate()
	utils.Print(paginated)

	updateUserById := repository.updateById()
	utils.Print(updateUserById)

	updatedUsers := repository.updateAll()
	utils.Print(updatedUsers)

	deletedUserById := repository.deleteById()
	utils.Print(deletedUserById)

	deletedUsers := repository.deleteAll()
	utils.Print(deletedUsers)
}

package goqu

import (
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/ofadiman/pg-go/seed"
	"github.com/ofadiman/pg-go/src/utils"
)

func UpdateUserById() {
	seed.Seed()

	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	db := goqu.New("postgres", sqlxdb.DB)

	var user User
	statement := db.Update("users").
		Set(
			UpdateUserData{
				FirstName: "Jack",
				LastName:  "Sparrow",
				UpdatedAt: time.Now(),
			},
		).
		Where(goqu.C("id").Eq("4b821cd3-47a0-44b8-9dbe-fe3a7625a163")).
		Returning(goqu.Star())

	_, err = statement.Executor().ScanStruct(&user)
	if err != nil {
		panic(err)
	}

	utils.Save(utils.GetPackageName(), user)
}

func UpdateUnderageUsers() {
	seed.Seed()

	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	db := goqu.New("postgres", sqlxdb.DB)

	var users []User
	statement := db.Update("users").
		Set(
			UpdateUnderageUsersData{
				EmailVerified: false,
				UpdatedAt:     time.Now(),
			},
		).
		Where(goqu.Literal("age(created_at)").Lte(goqu.Literal("interval '18 years'"))).
		Returning(goqu.Star())

	err = statement.Executor().ScanStructs(&users)
	if err != nil {
		panic(err)
	}

	utils.Save(utils.GetPackageName(), users)
}

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
			UpdateUser{
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

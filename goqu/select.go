package goqu

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/ofadiman/pg-go/seed"
	"github.com/ofadiman/pg-go/src/utils"
)

func SelectUserByID() {
	seed.Seed()

	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	db := goqu.New("postgres", sqlxdb.DB)

	var user User
	statement := db.From("users").
		Where(goqu.C("id").Eq("a5aef435-0ac5-4093-bcc5-71334bd8ffe8"))

	_, err = statement.Executor().ScanStruct(&user)
	if err != nil {
		panic(err)
	}

	utils.Save(utils.GetPackageName(), user)
}

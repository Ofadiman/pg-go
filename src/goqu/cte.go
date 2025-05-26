package main

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ofadiman/pg-go/src/utils"
)

func Cte() {
	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	database := goqu.New("postgres", sqlxdb.DB)

	statement := database.From("posts").
		With("oldest_users", goqu.From("users").Select(goqu.I("id")).Order(goqu.I("users.birthday").Asc()).Limit(3)).
		Select(goqu.Star()).Where(goqu.Ex{
		"author_id": goqu.From("oldest_users").Select("id"),
	})

	var posts []Post
	err = statement.Executor().ScanStructs(&posts)
	if err != nil {
		panic(err)
	}
	utils.Print(posts)
}

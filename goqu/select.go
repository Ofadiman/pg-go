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
	statement := db.Select(goqu.Star()).
		From("users").
		Where(goqu.C("id").Eq("a5aef435-0ac5-4093-bcc5-71334bd8ffe8"))

	_, err = statement.Executor().ScanStruct(&user)
	if err != nil {
		panic(err)
	}

	utils.Save(utils.GetPackageName(), user)
}

func SelectUsersFromOmahaBostonIndianapolis() {
	seed.Seed()

	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	db := goqu.New("postgres", sqlxdb.DB)

	var users []User
	statement := db.Select(goqu.Star()).
		From("users").
		Where(goqu.L("address ->> 'city'").In([]string{"Omaha", "Boston", "Indianapolis"}))

	err = statement.Executor().ScanStructs(&users)
	if err != nil {
		panic(err)
	}

	utils.Save(utils.GetPackageName(), users)
}

func SelectOldEnoughVerifiedMalesAndFemales() {
	seed.Seed()

	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	db := goqu.New("postgres", sqlxdb.DB)

	var users []User
	statement := db.Select(goqu.Star()).
		From("users").
		Where(goqu.C("email_verified").Eq(true)).
		Where(goqu.Or(
			goqu.And(
				goqu.C("gender").Eq("male"),
				goqu.L("age(birthday)").Gte(goqu.L("interval '21 years'")),
			),
			goqu.And(
				goqu.C("gender").Eq("female"),
				goqu.L("age(birthday)").Gte(goqu.L("interval '18 years'")),
			),
		))

	err = statement.Executor().ScanStructs(&users)
	if err != nil {
		panic(err)
	}

	utils.Save(utils.GetPackageName(), users)
}

func SelectAnStart() {
	seed.Seed()

	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	db := goqu.New("postgres", sqlxdb.DB)

	var users []User
	statement := db.Select(goqu.Star()).
		From("users").
		Where(goqu.C("first_name").Like("An%"))

	err = statement.Executor().ScanStructs(&users)
	if err != nil {
		panic(err)
	}

	utils.Save(utils.GetPackageName(), users)
}

func SelectNotOlderThanEighteenNorVerified() {
	seed.Seed()

	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	db := goqu.New("postgres", sqlxdb.DB)

	var users []User
	statement := db.Select(goqu.Star()).
		From("users").
		Where(goqu.L("age(birthday)").Lte(goqu.L("interval '18 years'"))).
		Where(goqu.C("email_verified").Neq(true))

	err = statement.Executor().ScanStructs(&users)
	if err != nil {
		panic(err)
	}

	utils.Save(utils.GetPackageName(), users)
}

func SelectUserIDsByCountry() {
	seed.Seed()

	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	db := goqu.New("postgres", sqlxdb.DB)

	var userIDsByCountry []UserIdsByCountry
	statement := db.Select(
		goqu.COUNT("*").As("total"),
		goqu.Literal("address ->> 'country'").As("country"),
		goqu.Literal("json_agg(id)").As("user_ids"),
	).
		From("users").
		GroupBy(goqu.L("address ->> 'country'")).
		Order(goqu.C("total").Desc())

	err = statement.Executor().ScanStructs(&userIDsByCountry)
	if err != nil {
		panic(err)
	}

	utils.Save(utils.GetPackageName(), userIDsByCountry)
}

func SelectThreeMostPopularAgeGroups() {
	seed.Seed()

	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	db := goqu.New("postgres", sqlxdb.DB)

	var mostPopularAgeGroups []MostPopularAgeGroup
	statement := db.From("ranked_users").
		With(
			"ranked_users",
			goqu.From("users").
				Select(
					goqu.Literal("extract('year' from age(birthday))").As("age"),
					goqu.COUNT("*").As("total"),
					goqu.Literal("dense_rank() OVER (ORDER BY count(*) DESC)").As("rank"),
				).
				GroupBy(goqu.Literal("extract('year' from age(birthday))")),
		).
		Select(
			goqu.C("age"),
			goqu.C("total"),
		).
		Order(goqu.C("total").Desc(), goqu.C("age").Desc()).
		Where(goqu.C("rank").Lte(3))

	err = statement.Executor().ScanStructs(&mostPopularAgeGroups)
	if err != nil {
		panic(err)
	}

	utils.Save(utils.GetPackageName(), mostPopularAgeGroups)
}

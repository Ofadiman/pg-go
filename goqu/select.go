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

func SelectMostRecentPostsWithRelatedRecords() {
	seed.Seed()

	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	db := goqu.New("postgres", sqlxdb.DB)

	var recentPostsWithRelatedRecords []RecentPostWithRelatedRecords
	statement := db.From("recent_posts").
		With(
			"recent_posts",
			goqu.From("posts").
				Select(
					goqu.Star(),
				).
				Order(goqu.C("created_at").Desc()).
				Limit(10),
		).
		Select(
			goqu.T("recent_posts").Col("id").As("post_id"),
			goqu.T("recent_posts").Col("user_id").As("post_user_id"),
			goqu.T("recent_posts").Col("title").As("post_title"),
			goqu.T("recent_posts").Col("content").As("post_content"),
			goqu.T("recent_posts").Col("created_at").As("post_created_at"),
			goqu.T("recent_posts").Col("updated_at").As("post_updated_at"),
			goqu.T("users").Col("id").As("user_id"),
			goqu.T("users").Col("created_at").As("user_created_at"),
			goqu.T("users").Col("updated_at").As("user_updated_at"),
			goqu.T("users").Col("address").As("user_address"),
			goqu.T("users").Col("birthday").As("user_birthday"),
			goqu.T("users").Col("email").As("user_email"),
			goqu.T("users").Col("email_verified").As("user_email_verified"),
			goqu.T("users").Col("first_name").As("user_first_name"),
			goqu.T("users").Col("last_name").As("user_last_name"),
			goqu.T("users").Col("gender").As("user_gender"),
			goqu.T("comments").Col("id").As("comment_id"),
			goqu.T("comments").Col("user_id").As("comment_user_id"),
			goqu.T("comments").Col("parent_id").As("comment_parent_id"),
			goqu.T("comments").Col("post_id").As("comment_post_id"),
			goqu.T("comments").Col("created_at").As("comment_created_at"),
			goqu.T("comments").Col("updated_at").As("comment_updated_at"),
			goqu.T("comments").Col("content").As("comment_content"),
			goqu.T("reactions").Col("id").As("reaction_id"),
			goqu.T("reactions").Col("comment_id").As("reaction_comment_id"),
			goqu.T("reactions").Col("user_id").As("reaction_user_id"),
			goqu.T("reactions").Col("created_at").As("reaction_created_at"),
			goqu.T("reactions").Col("emoji").As("reaction_emoji"),
		).
		LeftJoin(
			goqu.T("users"),
			goqu.On(goqu.T("users").Col("id").Eq(goqu.T("recent_posts").Col("user_id"))),
		).
		LeftJoin(
			goqu.T("comments"),
			goqu.On(goqu.T("comments").Col("post_id").Eq(goqu.T("recent_posts").Col("id"))),
		).
		LeftJoin(
			goqu.T("reactions"),
			goqu.On(goqu.T("reactions").Col("comment_id").Eq(goqu.T("comments").Col("id"))),
		)

	err = statement.Executor().ScanStructs(&recentPostsWithRelatedRecords)
	if err != nil {
		panic(err)
	}

	utils.Save(utils.GetPackageName(), recentPostsWithRelatedRecords)
}

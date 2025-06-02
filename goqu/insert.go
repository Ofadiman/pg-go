package goqu

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/ofadiman/pg-go/src/utils"
)

func InsertUser() {
	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	sqlxdb.Exec("truncate table users cascade")

	db := goqu.New("postgres", sqlxdb.DB)

	var user User = User{
		FirstName:     gofakeit.FirstName(),
		LastName:      gofakeit.LastName(),
		Email:         gofakeit.Email(),
		EmailVerified: gofakeit.Bool(),
		Gender:        gofakeit.Gender(),
		Birthday:      gofakeit.PastDate(),
		CreatedAt:     gofakeit.PastDate(),
		Address: Address{
			City:         gofakeit.City(),
			Country:      gofakeit.Country(),
			Street:       gofakeit.Street(),
			StreetNumber: gofakeit.StreetNumber(),
			ZipCode:      gofakeit.Zip(),
		},
	}
	statement := db.Insert("users").Rows(user).Returning(goqu.Star())
	_, err = statement.Executor().ScanStruct(&user)
	if err != nil {
		panic(err)
	}

	utils.Save(utils.GetPackageName(), user)
}

func InsertManyUsers() {
	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	sqlxdb.Exec("truncate table users cascade")

	db := goqu.New("postgres", sqlxdb.DB)

	statement := db.Insert("users").Rows([]User{
		{
			FirstName:     gofakeit.FirstName(),
			LastName:      gofakeit.LastName(),
			Email:         gofakeit.Email(),
			EmailVerified: gofakeit.Bool(),
			Gender:        gofakeit.Gender(),
			Birthday:      gofakeit.PastDate(),
			CreatedAt:     gofakeit.PastDate(),
			Address: Address{
				City:         gofakeit.City(),
				Country:      gofakeit.Country(),
				Street:       gofakeit.Street(),
				StreetNumber: gofakeit.StreetNumber(),
				ZipCode:      gofakeit.Zip(),
			},
		},
		{
			FirstName:     gofakeit.FirstName(),
			LastName:      gofakeit.LastName(),
			Email:         gofakeit.Email(),
			EmailVerified: gofakeit.Bool(),
			Gender:        gofakeit.Gender(),
			Birthday:      gofakeit.PastDate(),
			CreatedAt:     gofakeit.PastDate(),
			Address: Address{
				City:         gofakeit.City(),
				Country:      gofakeit.Country(),
				Street:       gofakeit.Street(),
				StreetNumber: gofakeit.StreetNumber(),
				ZipCode:      gofakeit.Zip(),
			},
		},
		{
			FirstName:     gofakeit.FirstName(),
			LastName:      gofakeit.LastName(),
			Email:         gofakeit.Email(),
			EmailVerified: gofakeit.Bool(),
			Gender:        gofakeit.Gender(),
			Birthday:      gofakeit.PastDate(),
			CreatedAt:     gofakeit.PastDate(),
			Address: Address{
				City:         gofakeit.City(),
				Country:      gofakeit.Country(),
				Street:       gofakeit.Street(),
				StreetNumber: gofakeit.StreetNumber(),
				ZipCode:      gofakeit.Zip(),
			},
		},
	}).
		Returning(goqu.Star())

	var users []User
	err = statement.Executor().ScanStructs(&users)
	if err != nil {
		panic(err)
	}

	utils.Save(utils.GetPackageName(), users)
}

func InsertUserWithPosts() {
	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	sqlxdb.Exec("truncate table users cascade")

	db := goqu.New("postgres", sqlxdb.DB)

	tx, err := db.Begin()
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	insertUserStatement := tx.Insert("users").Rows(User{
		FirstName:     gofakeit.FirstName(),
		LastName:      gofakeit.LastName(),
		Email:         gofakeit.Email(),
		EmailVerified: gofakeit.Bool(),
		Gender:        gofakeit.Gender(),
		Birthday:      gofakeit.PastDate(),
		Address: Address{
			City:         gofakeit.City(),
			Country:      gofakeit.Country(),
			Street:       gofakeit.Street(),
			StreetNumber: gofakeit.StreetNumber(),
			ZipCode:      gofakeit.Zip(),
		},
	}).
		Returning(goqu.Star())

	var user User
	_, err = insertUserStatement.Executor().ScanStruct(&user)
	if err != nil {
		panic(err)
	}

	insertPostsStatement := tx.Insert("posts").Rows([]Post{
		{
			UserID:  user.ID,
			Title:   gofakeit.LoremIpsumSentence(5),
			Content: gofakeit.LoremIpsumSentence(10),
		},
		{
			UserID:  user.ID,
			Title:   gofakeit.LoremIpsumSentence(5),
			Content: gofakeit.LoremIpsumSentence(10),
		},
		{
			UserID:  user.ID,
			Title:   gofakeit.LoremIpsumSentence(5),
			Content: gofakeit.LoremIpsumSentence(10),
		},
	}).Returning(goqu.Star())

	var posts []Post
	err = insertPostsStatement.Executor().ScanStructs(&posts)
	if err != nil {
		panic(err)
	}

	utils.Save(utils.GetPackageName(), []any{user, posts})
}

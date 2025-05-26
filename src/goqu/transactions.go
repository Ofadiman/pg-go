package main

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Transactions() {
	sqlxdb, err := sqlx.Connect("postgres", "postgres://admin:password@localhost:5432/pg-go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlxdb.Close()

	sqlxdb.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")

	database := goqu.New("postgres", sqlxdb.DB)

	tx, err := database.Begin()
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

	var user User
	_, err = tx.Insert("users").Rows(
		User{
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
		},
	).Returning(goqu.Star()).
		Executor().
		ScanStruct(&user)
	if err != nil {
		panic(err)
	}

	var posts []Post
	err = tx.Insert("posts").Rows([]Post{
		{
			ID:        gofakeit.UUID(),
			AuthorID:  user.ID,
			CreatedAt: gofakeit.PastDate(),
			Title:     gofakeit.LoremIpsumSentence(gofakeit.IntN(9) + 1),
			Content:   gofakeit.LoremIpsumParagraph(1, gofakeit.IntN(4)+1, gofakeit.IntN(9)+1, " "),
		},
		{
			ID:        gofakeit.UUID(),
			AuthorID:  user.ID,
			CreatedAt: gofakeit.PastDate(),
			Title:     gofakeit.LoremIpsumSentence(gofakeit.IntN(9) + 1),
			Content:   gofakeit.LoremIpsumParagraph(1, gofakeit.IntN(4)+1, gofakeit.IntN(9)+1, " "),
		},
		{
			ID:        gofakeit.UUID(),
			AuthorID:  user.ID,
			CreatedAt: gofakeit.PastDate(),
			Title:     gofakeit.LoremIpsumSentence(gofakeit.IntN(9) + 1),
			Content:   gofakeit.LoremIpsumParagraph(1, gofakeit.IntN(4)+1, gofakeit.IntN(9)+1, " "),
		},
	}).
		Returning(goqu.Star()).
		Executor().
		ScanStructs(&posts)
	if err != nil {
		panic(err)
	}
}

package main

import "time"

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

type Post struct {
	ID       string `db:"id"`
	AuthorID string `db:"author_id"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`

	Title   string `db:"title"`
	Content string `db:"content"`
}

func main() {
	Cte()
}

package integration

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ofadiman/pg-go/repository"
)

const (
	missingUserID    = "00000000-0000-0000-0000-999999999999"
	missingPostID    = "00000000-0000-0000-0001-999999999999"
	missingCommentID = "00000000-0000-0000-0002-999999999999"
)

func newUser(n int) repository.NewUser {
	gender := "male"
	if n%2 == 0 {
		gender = "female"
	}
	return repository.NewUser{
		Address: repository.Address{
			City:         "Gdansk",
			Country:      "Poland",
			Street:       "Long Market",
			StreetNumber: strconv.Itoa(n),
			ZipCode:      fmt.Sprintf("80-%03d", n),
		},
		Birthday:      time.Date(1990+n%20, time.Month(1+n%12), 1+n%28, 0, 0, 0, 0, time.UTC),
		Email:         fmt.Sprintf("new-user-%d@example.com", n),
		EmailVerified: n%2 == 0,
		FirstName:     fmt.Sprintf("First%d", n),
		Gender:        gender,
		LastName:      fmt.Sprintf("Last%d", n),
	}
}

func newPost(userID string, n int) repository.NewPost {
	return repository.NewPost{
		UserID:  userID,
		Title:   fmt.Sprintf("Transaction post %d", n),
		Content: fmt.Sprintf("Content of transaction post %d.", n),
	}
}

package goqu

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Address struct {
	City         string `json:"city"`
	Country      string `json:"country"`
	Street       string `json:"street"`
	StreetNumber string `json:"street_number"`
	ZipCode      string `json:"zip_code"`
}

func (a Address) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Address) Scan(value any) error {
	if value == nil {
		return nil
	}

	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return errors.New("cannot scan into Address: unsupported type")
	}

	return json.Unmarshal(b, a)
}

type User struct {
	ID string `db:"id"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`

	Address       Address   `db:"address"`
	Birthday      time.Time `db:"birthday"`
	Email         string    `db:"email"`
	EmailVerified bool      `db:"email_verified"`
	FirstName     string    `db:"first_name"`
	Gender        string    `db:"gender"`
	LastName      string    `db:"last_name"`
}

type UserIDs []string

func (u *UserIDs) Scan(value any) error {
	if value == nil {
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}

	return json.Unmarshal(data, u)
}

func (u UserIDs) Value() (driver.Value, error) {
	return json.Marshal(u)
}

type UserIdsByCountry struct {
	Country string  `db:"country"`
	Total   int     `db:"total"`
	UserIDs UserIDs `db:"user_ids"`
}

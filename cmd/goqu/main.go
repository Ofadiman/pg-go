package main

import (
	"github.com/ofadiman/pg-go/goqu"
)

func main() {
	goqu.SelectUserByID()
	goqu.SelectUsersFromOmahaBostonIndianapolis()
	goqu.SelectOldEnoughVerifiedMalesAndFemales()
}

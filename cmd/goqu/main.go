package main

import (
	"github.com/ofadiman/pg-go/goqu"
)

func main() {
	goqu.SelectUserByID()
	goqu.SelectUsersFromOmahaBostonIndianapolis()
	goqu.SelectOldEnoughVerifiedMalesAndFemales()
	goqu.SelectAnStart()
	goqu.SelectNotOlderThanEighteenNorVerified()
	goqu.SelectUserIDsByCountry()
	goqu.SelectThreeMostPopularAgeGroups()
}

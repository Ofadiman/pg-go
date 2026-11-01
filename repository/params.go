package repository

import (
	"fmt"
	"time"
)

const (
	defaultLimit = 20
	maxLimit     = 100
	maxSorts     = 2
)

type ListUsersParams struct {
	FirstNamePrefix *string
	Gender          *string
	EmailVerified   *bool
	BornAfter       *time.Time
	BornBefore      *time.Time
	Sort            []UserSort
	Limit           int
	Offset          int
}

type UserSort struct {
	Field UserSortField
	Desc  bool
}

type UserSortField string

const (
	SortByFirstName UserSortField = "first_name"
	SortByLastName  UserSortField = "last_name"
	SortByBirthday  UserSortField = "birthday"
	SortByCreatedAt UserSortField = "created_at"
)

type Page[T any] struct {
	Items []T
	Total int
}

func (p ListUsersParams) Normalize() (ListUsersParams, error) {
	if p.Limit == 0 {
		p.Limit = defaultLimit
	}
	if p.Limit < 1 || p.Limit > maxLimit {
		return ListUsersParams{}, fmt.Errorf("%w: limit must be between 1 and %d, got %d", ErrInvalidParams, maxLimit, p.Limit)
	}
	if p.Offset < 0 {
		return ListUsersParams{}, fmt.Errorf("%w: offset must not be negative, got %d", ErrInvalidParams, p.Offset)
	}
	if len(p.Sort) > maxSorts {
		return ListUsersParams{}, fmt.Errorf("%w: at most %d sort fields are allowed, got %d", ErrInvalidParams, maxSorts, len(p.Sort))
	}

	seen := make(map[UserSortField]bool, len(p.Sort))
	for _, sort := range p.Sort {
		switch sort.Field {
		case SortByFirstName, SortByLastName, SortByBirthday, SortByCreatedAt:
		default:
			return ListUsersParams{}, fmt.Errorf("%w: unknown sort field %q", ErrInvalidParams, sort.Field)
		}
		if seen[sort.Field] {
			return ListUsersParams{}, fmt.Errorf("%w: duplicate sort field %q", ErrInvalidParams, sort.Field)
		}
		seen[sort.Field] = true
	}

	if len(p.Sort) == 0 {
		p.Sort = []UserSort{{Field: SortByCreatedAt, Desc: true}}
	}

	return p, nil
}

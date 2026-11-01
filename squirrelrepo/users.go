package squirrelrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/ofadiman/pg-go/repository"
)

type users struct {
	*Store
}

func (r users) insert(in []repository.NewUser, now time.Time) squirrel.InsertBuilder {
	query := r.sb.Insert("users").
		Columns("address", "birthday", "created_at", "email", "email_verified", "first_name", "gender", "last_name")
	for _, user := range in {
		query = query.Values(user.Address, user.Birthday, now, user.Email, user.EmailVerified, user.FirstName, user.Gender, user.LastName)
	}
	return query
}

func (r users) Create(ctx context.Context, in repository.NewUser) (repository.User, error) {
	query := r.insert([]repository.NewUser{in}, r.now()).Suffix("RETURNING *")

	var row userRow
	if err := r.get(ctx, &row, query); err != nil {
		return repository.User{}, fmt.Errorf("create user: %w", translate(err))
	}
	return row.toDomain(), nil
}

func (r users) CreateMany(ctx context.Context, in []repository.NewUser) ([]repository.User, error) {
	query := r.insert(in, r.now()).Suffix("RETURNING *")

	var rows []userRow
	if err := r.selectAll(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("create users: %w", translate(err))
	}
	return usersToDomain(rows), nil
}

func (r users) GetByID(ctx context.Context, id string) (repository.User, error) {
	query := r.sb.Select(userColumns...).
		From("users").
		Where(squirrel.Eq{"id": id})

	var row userRow
	if err := r.get(ctx, &row, query); err != nil {
		if notFound(err) {
			return repository.User{}, fmt.Errorf("get user %s: %w", id, repository.ErrNotFound)
		}
		return repository.User{}, fmt.Errorf("get user %s: %w", id, err)
	}
	return row.toDomain(), nil
}

func (r users) Update(ctx context.Context, id string, in repository.UserUpdate) (repository.User, error) {
	query := r.sb.Update("users").Set("updated_at", r.now())
	if in.FirstName != nil {
		query = query.Set("first_name", *in.FirstName)
	}
	if in.LastName != nil {
		query = query.Set("last_name", *in.LastName)
	}
	if in.EmailVerified != nil {
		query = query.Set("email_verified", *in.EmailVerified)
	}
	query = query.Where(squirrel.Eq{"id": id}).Suffix("RETURNING *")

	var row userRow
	if err := r.get(ctx, &row, query); err != nil {
		if notFound(err) {
			return repository.User{}, fmt.Errorf("update user %s: %w", id, repository.ErrNotFound)
		}
		return repository.User{}, fmt.Errorf("update user %s: %w", id, err)
	}
	return row.toDomain(), nil
}

func (r users) Delete(ctx context.Context, id string) (repository.User, error) {
	query := r.sb.Delete("users").
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING *")

	var row userRow
	if err := r.get(ctx, &row, query); err != nil {
		if notFound(err) {
			return repository.User{}, fmt.Errorf("delete user %s: %w", id, repository.ErrNotFound)
		}
		return repository.User{}, fmt.Errorf("delete user %s: %w", id, err)
	}
	return row.toDomain(), nil
}

func (r users) Upsert(ctx context.Context, in repository.NewUser) (repository.User, error) {
	now := r.now()
	query := r.insert([]repository.NewUser{in}, now).
		Suffix("ON CONFLICT (email) DO UPDATE SET first_name = excluded.first_name, last_name = excluded.last_name, updated_at = ? RETURNING *", now)

	var row userRow
	if err := r.get(ctx, &row, query); err != nil {
		return repository.User{}, fmt.Errorf("upsert user: %w", err)
	}
	return row.toDomain(), nil
}

func (r users) List(ctx context.Context, params repository.ListUsersParams) (repository.Page[repository.User], error) {
	params, err := params.Normalize()
	if err != nil {
		return repository.Page[repository.User]{}, fmt.Errorf("list users: %w", err)
	}

	base := r.sb.Select().From("users")
	if params.FirstNamePrefix != nil {
		base = base.Where(squirrel.Like{"first_name": *params.FirstNamePrefix + "%"})
	}
	if params.Gender != nil {
		base = base.Where(squirrel.Eq{"gender": *params.Gender})
	}
	if params.EmailVerified != nil {
		base = base.Where(squirrel.Eq{"email_verified": *params.EmailVerified})
	}
	if params.BornAfter != nil {
		base = base.Where(squirrel.Gt{"birthday": *params.BornAfter})
	}
	if params.BornBefore != nil {
		base = base.Where(squirrel.Lt{"birthday": *params.BornBefore})
	}

	var total int
	if err := r.get(ctx, &total, base.Columns("count(*)")); err != nil {
		return repository.Page[repository.User]{}, fmt.Errorf("count users: %w", err)
	}

	order := make([]string, 0, len(params.Sort)+1)
	for _, sort := range params.Sort {
		direction := "ASC"
		if sort.Desc {
			direction = "DESC"
		}
		order = append(order, string(sort.Field)+" "+direction)
	}
	order = append(order, "id ASC")

	query := base.Columns(userColumns...).
		OrderBy(order...).
		Limit(uint64(params.Limit)).
		Offset(uint64(params.Offset))

	var rows []userRow
	if err := r.selectAll(ctx, &rows, query); err != nil {
		return repository.Page[repository.User]{}, fmt.Errorf("list users: %w", err)
	}

	return repository.Page[repository.User]{Items: usersToDomain(rows), Total: total}, nil
}

func (r users) ListEligibleAdults(ctx context.Context) ([]repository.User, error) {
	now := r.now()
	query := r.sb.Select(userColumns...).
		From("users").
		Where(squirrel.Eq{"email_verified": true}).
		Where(squirrel.Or{
			squirrel.And{
				squirrel.Eq{"gender": "male"},
				squirrel.Expr("age(?, birthday) >= interval '21 years'", now),
			},
			squirrel.And{
				squirrel.Eq{"gender": "female"},
				squirrel.Expr("age(?, birthday) >= interval '18 years'", now),
			},
		}).
		OrderBy("id ASC")

	var rows []userRow
	if err := r.selectAll(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("list eligible adults: %w", err)
	}
	return usersToDomain(rows), nil
}

func (r users) UnverifyMinors(ctx context.Context) ([]repository.User, error) {
	now := r.now()
	query := r.sb.Update("users").
		Set("email_verified", false).
		Set("updated_at", now).
		Where("age(?, birthday) < interval '18 years'", now).
		Suffix("RETURNING *")

	var rows []userRow
	if err := r.selectAll(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("unverify minors: %w", err)
	}
	return usersToDomain(rows), nil
}

func (r users) TopAgeGroups(ctx context.Context, limit int) ([]repository.AgeGroup, error) {
	ranked := squirrel.Select().
		Column(squirrel.Expr("extract(year from age(?, birthday))::int AS age", r.now())).
		Column("count(*) AS total").
		Column("dense_rank() over (order by count(*) desc) AS rank").
		From("users").
		GroupBy("age")

	query := r.sb.Select("age", "total").
		PrefixExpr(squirrel.Expr("WITH ranked AS (?)", ranked)).
		From("ranked").
		Where(squirrel.LtOrEq{"rank": limit}).
		OrderBy("total DESC", "age ASC")

	var rows []ageGroupRow
	if err := r.selectAll(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("top age groups: %w", err)
	}

	out := make([]repository.AgeGroup, len(rows))
	for i, row := range rows {
		out[i] = repository.AgeGroup{Age: row.Age, Total: row.Total}
	}
	return out, nil
}

func (r users) PostStats(ctx context.Context, minPosts int) ([]repository.UserPostStats, error) {
	query := r.sb.Select(
		"u.id AS user_id",
		"u.first_name",
		"u.last_name",
		"count(distinct p.id) AS post_count",
		"count(distinct c.id) AS comment_count",
		"max(p.created_at) AS last_post_at",
	).
		From("users AS u").
		Join("posts AS p ON p.user_id = u.id").
		LeftJoin("comments AS c ON c.post_id = p.id").
		GroupBy("u.id").
		Having("count(distinct p.id) >= ?", minPosts).
		OrderBy("post_count DESC", "u.id ASC")

	var rows []userPostStatsRow
	if err := r.selectAll(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("post stats: %w", err)
	}

	out := make([]repository.UserPostStats, len(rows))
	for i, row := range rows {
		out[i] = repository.UserPostStats(row)
	}
	return out, nil
}

func (r users) ClaimUnverified(ctx context.Context, limit int) ([]repository.User, error) {
	query := r.sb.Select(userColumns...).
		From("users").
		Where(squirrel.Eq{"email_verified": false}).
		OrderBy("created_at ASC", "id ASC").
		Limit(uint64(limit)).
		Suffix("FOR UPDATE SKIP LOCKED")

	var rows []userRow
	if err := r.selectAll(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("claim unverified users: %w", err)
	}
	return usersToDomain(rows), nil
}

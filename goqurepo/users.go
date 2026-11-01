package goqurepo

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/ofadiman/pg-go/repository"
)

type users struct {
	*Store
}

func (r users) Create(ctx context.Context, in repository.NewUser) (repository.User, error) {
	var row userRow
	_, err := r.q.Insert("users").
		Rows(newUserRow(in, r.now())).
		Returning(goqu.Star()).
		Executor().
		ScanStructContext(ctx, &row)
	if err != nil {
		return repository.User{}, fmt.Errorf("create user: %w", translate(err))
	}
	return row.toDomain(), nil
}

func (r users) CreateMany(ctx context.Context, in []repository.NewUser) ([]repository.User, error) {
	now := r.now()
	rows := make([]userRow, len(in))
	for i, user := range in {
		rows[i] = newUserRow(user, now)
	}

	var created []userRow
	err := r.q.Insert("users").
		Rows(rows).
		Returning(goqu.Star()).
		Executor().
		ScanStructsContext(ctx, &created)
	if err != nil {
		return nil, fmt.Errorf("create users: %w", translate(err))
	}
	return usersToDomain(created), nil
}

func (r users) GetByID(ctx context.Context, id string) (repository.User, error) {
	var row userRow
	found, err := r.q.From("users").
		Where(goqu.C("id").Eq(id)).
		ScanStructContext(ctx, &row)
	if err != nil {
		return repository.User{}, fmt.Errorf("get user %s: %w", id, err)
	}
	if !found {
		return repository.User{}, fmt.Errorf("get user %s: %w", id, repository.ErrNotFound)
	}
	return row.toDomain(), nil
}

func (r users) Update(ctx context.Context, id string, in repository.UserUpdate) (repository.User, error) {
	record := goqu.Record{"updated_at": r.now()}
	if in.FirstName != nil {
		record["first_name"] = *in.FirstName
	}
	if in.LastName != nil {
		record["last_name"] = *in.LastName
	}
	if in.EmailVerified != nil {
		record["email_verified"] = *in.EmailVerified
	}

	var row userRow
	found, err := r.q.Update("users").
		Set(record).
		Where(goqu.C("id").Eq(id)).
		Returning(goqu.Star()).
		Executor().
		ScanStructContext(ctx, &row)
	if err != nil {
		return repository.User{}, fmt.Errorf("update user %s: %w", id, err)
	}
	if !found {
		return repository.User{}, fmt.Errorf("update user %s: %w", id, repository.ErrNotFound)
	}
	return row.toDomain(), nil
}

func (r users) Delete(ctx context.Context, id string) (repository.User, error) {
	var row userRow
	found, err := r.q.Delete("users").
		Where(goqu.C("id").Eq(id)).
		Returning(goqu.Star()).
		Executor().
		ScanStructContext(ctx, &row)
	if err != nil {
		return repository.User{}, fmt.Errorf("delete user %s: %w", id, err)
	}
	if !found {
		return repository.User{}, fmt.Errorf("delete user %s: %w", id, repository.ErrNotFound)
	}
	return row.toDomain(), nil
}

func (r users) Upsert(ctx context.Context, in repository.NewUser) (repository.User, error) {
	now := r.now()

	var row userRow
	_, err := r.q.Insert("users").
		Rows(newUserRow(in, now)).
		OnConflict(goqu.DoUpdate("email", goqu.Record{
			"first_name": goqu.I("excluded.first_name"),
			"last_name":  goqu.I("excluded.last_name"),
			"updated_at": now,
		})).
		Returning(goqu.Star()).
		Executor().
		ScanStructContext(ctx, &row)
	if err != nil {
		return repository.User{}, fmt.Errorf("upsert user: %w", err)
	}
	return row.toDomain(), nil
}

func (r users) List(ctx context.Context, params repository.ListUsersParams) (repository.Page[repository.User], error) {
	params, err := params.Normalize()
	if err != nil {
		return repository.Page[repository.User]{}, fmt.Errorf("list users: %w", err)
	}

	var filters []exp.Expression
	if params.FirstNamePrefix != nil {
		filters = append(filters, goqu.C("first_name").Like(*params.FirstNamePrefix+"%"))
	}
	if params.Gender != nil {
		filters = append(filters, goqu.C("gender").Eq(*params.Gender))
	}
	if params.EmailVerified != nil {
		filters = append(filters, goqu.C("email_verified").Eq(*params.EmailVerified))
	}
	if params.BornAfter != nil {
		filters = append(filters, goqu.C("birthday").Gt(*params.BornAfter))
	}
	if params.BornBefore != nil {
		filters = append(filters, goqu.C("birthday").Lt(*params.BornBefore))
	}

	base := r.q.From("users").Where(filters...)

	total, err := base.CountContext(ctx)
	if err != nil {
		return repository.Page[repository.User]{}, fmt.Errorf("count users: %w", err)
	}

	order := make([]exp.OrderedExpression, 0, len(params.Sort)+1)
	for _, sort := range params.Sort {
		column := goqu.C(string(sort.Field))
		if sort.Desc {
			order = append(order, column.Desc())
		} else {
			order = append(order, column.Asc())
		}
	}
	order = append(order, goqu.C("id").Asc())

	var rows []userRow
	err = base.
		Order(order...).
		Limit(uint(params.Limit)).
		Offset(uint(params.Offset)).
		ScanStructsContext(ctx, &rows)
	if err != nil {
		return repository.Page[repository.User]{}, fmt.Errorf("list users: %w", err)
	}

	return repository.Page[repository.User]{Items: usersToDomain(rows), Total: int(total)}, nil
}

func (r users) ListEligibleAdults(ctx context.Context) ([]repository.User, error) {
	age := goqu.L("age(?, birthday)", r.now())

	var rows []userRow
	err := r.q.From("users").
		Where(
			goqu.C("email_verified").Eq(true),
			goqu.Or(
				goqu.And(goqu.C("gender").Eq("male"), age.Gte(goqu.L("interval '21 years'"))),
				goqu.And(goqu.C("gender").Eq("female"), age.Gte(goqu.L("interval '18 years'"))),
			),
		).
		Order(goqu.C("id").Asc()).
		ScanStructsContext(ctx, &rows)
	if err != nil {
		return nil, fmt.Errorf("list eligible adults: %w", err)
	}
	return usersToDomain(rows), nil
}

func (r users) UnverifyMinors(ctx context.Context) ([]repository.User, error) {
	now := r.now()

	var rows []userRow
	err := r.q.Update("users").
		Set(goqu.Record{"email_verified": false, "updated_at": now}).
		Where(goqu.L("age(?, birthday)", now).Lt(goqu.L("interval '18 years'"))).
		Returning(goqu.Star()).
		Executor().
		ScanStructsContext(ctx, &rows)
	if err != nil {
		return nil, fmt.Errorf("unverify minors: %w", err)
	}
	return usersToDomain(rows), nil
}

func (r users) TopAgeGroups(ctx context.Context, limit int) ([]repository.AgeGroup, error) {
	ranked := r.q.From("users").
		Select(
			goqu.L("extract(year from age(?, birthday))::int", r.now()).As("age"),
			goqu.COUNT(goqu.Star()).As("total"),
			goqu.L("dense_rank() over (order by count(*) desc)").As("rank"),
		).
		GroupBy(goqu.C("age"))

	var rows []ageGroupRow
	err := r.q.From("ranked").
		With("ranked", ranked).
		Select(goqu.C("age"), goqu.C("total")).
		Where(goqu.C("rank").Lte(limit)).
		Order(goqu.C("total").Desc(), goqu.C("age").Asc()).
		ScanStructsContext(ctx, &rows)
	if err != nil {
		return nil, fmt.Errorf("top age groups: %w", err)
	}

	out := make([]repository.AgeGroup, len(rows))
	for i, row := range rows {
		out[i] = repository.AgeGroup{Age: row.Age, Total: row.Total}
	}
	return out, nil
}

func (r users) PostStats(ctx context.Context, minPosts int) ([]repository.UserPostStats, error) {
	postCount := goqu.COUNT(goqu.DISTINCT(goqu.I("p.id")))

	var rows []userPostStatsRow
	err := r.q.From(goqu.T("users").As("u")).
		Join(goqu.T("posts").As("p"), goqu.On(goqu.I("p.user_id").Eq(goqu.I("u.id")))).
		LeftJoin(goqu.T("comments").As("c"), goqu.On(goqu.I("c.post_id").Eq(goqu.I("p.id")))).
		Select(
			goqu.I("u.id").As("user_id"),
			goqu.I("u.first_name"),
			goqu.I("u.last_name"),
			postCount.As("post_count"),
			goqu.COUNT(goqu.DISTINCT(goqu.I("c.id"))).As("comment_count"),
			goqu.MAX(goqu.I("p.created_at")).As("last_post_at"),
		).
		GroupBy(goqu.I("u.id")).
		Having(postCount.Gte(minPosts)).
		Order(goqu.C("post_count").Desc(), goqu.I("u.id").Asc()).
		ScanStructsContext(ctx, &rows)
	if err != nil {
		return nil, fmt.Errorf("post stats: %w", err)
	}

	out := make([]repository.UserPostStats, len(rows))
	for i, row := range rows {
		out[i] = repository.UserPostStats(row)
	}
	return out, nil
}

func (r users) ClaimUnverified(ctx context.Context, limit int) ([]repository.User, error) {
	var rows []userRow
	err := r.q.From("users").
		Where(goqu.C("email_verified").Eq(false)).
		Order(goqu.C("created_at").Asc(), goqu.C("id").Asc()).
		Limit(uint(limit)).
		ForUpdate(exp.SkipLocked).
		ScanStructsContext(ctx, &rows)
	if err != nil {
		return nil, fmt.Errorf("claim unverified users: %w", err)
	}
	return usersToDomain(rows), nil
}

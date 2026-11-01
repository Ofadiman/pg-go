package gormrepo

import (
	"context"
	"fmt"

	"github.com/ofadiman/pg-go/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type users struct {
	*Store
}

func (r users) Create(ctx context.Context, in repository.NewUser) (repository.User, error) {
	row := newUserRow(in, r.now())
	if err := r.q.WithContext(ctx).Create(&row).Error; err != nil {
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

	if err := r.q.WithContext(ctx).Create(&rows).Error; err != nil {
		return nil, fmt.Errorf("create users: %w", translate(err))
	}
	return usersToDomain(rows), nil
}

func (r users) GetByID(ctx context.Context, id string) (repository.User, error) {
	var row userRow
	if err := r.q.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		if notFound(err) {
			return repository.User{}, fmt.Errorf("get user %s: %w", id, repository.ErrNotFound)
		}
		return repository.User{}, fmt.Errorf("get user %s: %w", id, err)
	}
	return row.toDomain(), nil
}

func (r users) Update(ctx context.Context, id string, in repository.UserUpdate) (repository.User, error) {
	values := map[string]any{"updated_at": r.now()}
	if in.FirstName != nil {
		values["first_name"] = *in.FirstName
	}
	if in.LastName != nil {
		values["last_name"] = *in.LastName
	}
	if in.EmailVerified != nil {
		values["email_verified"] = *in.EmailVerified
	}

	var row userRow
	result := r.q.WithContext(ctx).
		Model(&row).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(values)
	if result.Error != nil {
		return repository.User{}, fmt.Errorf("update user %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.User{}, fmt.Errorf("update user %s: %w", id, repository.ErrNotFound)
	}
	return row.toDomain(), nil
}

func (r users) Delete(ctx context.Context, id string) (repository.User, error) {
	var row userRow
	result := r.q.WithContext(ctx).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Delete(&row)
	if result.Error != nil {
		return repository.User{}, fmt.Errorf("delete user %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.User{}, fmt.Errorf("delete user %s: %w", id, repository.ErrNotFound)
	}
	return row.toDomain(), nil
}

func (r users) Upsert(ctx context.Context, in repository.NewUser) (repository.User, error) {
	now := r.now()
	row := newUserRow(in, now)
	err := r.q.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "email"}},
				DoUpdates: clause.Assignments(map[string]any{
					"first_name": gorm.Expr("excluded.first_name"),
					"last_name":  gorm.Expr("excluded.last_name"),
					"updated_at": now,
				}),
			},
			clause.Returning{},
		).
		Create(&row).Error
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

	base := r.q.WithContext(ctx).Model(&userRow{})
	if params.FirstNamePrefix != nil {
		base = base.Where("first_name LIKE ?", *params.FirstNamePrefix+"%")
	}
	if params.Gender != nil {
		base = base.Where("gender = ?", *params.Gender)
	}
	if params.EmailVerified != nil {
		base = base.Where("email_verified = ?", *params.EmailVerified)
	}
	if params.BornAfter != nil {
		base = base.Where("birthday > ?", *params.BornAfter)
	}
	if params.BornBefore != nil {
		base = base.Where("birthday < ?", *params.BornBefore)
	}
	base = base.Session(&gorm.Session{})

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return repository.Page[repository.User]{}, fmt.Errorf("count users: %w", err)
	}

	query := base
	for _, sort := range params.Sort {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: string(sort.Field)}, Desc: sort.Desc})
	}

	var rows []userRow
	err = query.
		Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}}).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&rows).Error
	if err != nil {
		return repository.Page[repository.User]{}, fmt.Errorf("list users: %w", err)
	}

	return repository.Page[repository.User]{Items: usersToDomain(rows), Total: int(total)}, nil
}

func (r users) ListEligibleAdults(ctx context.Context) ([]repository.User, error) {
	now := r.now()
	adults := r.q.
		Where("gender = ? AND age(?, birthday) >= interval '21 years'", "male", now).
		Or("gender = ? AND age(?, birthday) >= interval '18 years'", "female", now)

	var rows []userRow
	err := r.q.WithContext(ctx).
		Where("email_verified = ?", true).
		Where(adults).
		Order("id").
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list eligible adults: %w", err)
	}
	return usersToDomain(rows), nil
}

func (r users) UnverifyMinors(ctx context.Context) ([]repository.User, error) {
	now := r.now()

	var rows []userRow
	err := r.q.WithContext(ctx).
		Model(&rows).
		Clauses(clause.Returning{}).
		Where("age(?, birthday) < interval '18 years'", now).
		Updates(map[string]any{"email_verified": false, "updated_at": now}).Error
	if err != nil {
		return nil, fmt.Errorf("unverify minors: %w", err)
	}
	return usersToDomain(rows), nil
}

func (r users) TopAgeGroups(ctx context.Context, limit int) ([]repository.AgeGroup, error) {
	ranked := r.q.Model(&userRow{}).
		Select(
			"extract(year from age(?, birthday))::int AS age, count(*) AS total, dense_rank() over (order by count(*) desc) AS rank",
			r.now(),
		).
		Group("age")

	var rows []ageGroupRow
	err := r.q.WithContext(ctx).
		Raw("WITH ranked AS (?) SELECT age, total FROM ranked WHERE rank <= ? ORDER BY total DESC, age ASC", ranked, limit).
		Scan(&rows).Error
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
	var rows []userPostStatsRow
	err := r.q.WithContext(ctx).
		Table("users AS u").
		Select("u.id AS user_id, u.first_name, u.last_name, count(distinct p.id) AS post_count, count(distinct c.id) AS comment_count, max(p.created_at) AS last_post_at").
		Joins("JOIN posts AS p ON p.user_id = u.id").
		Joins("LEFT JOIN comments AS c ON c.post_id = p.id").
		Group("u.id").
		Having("count(distinct p.id) >= ?", minPosts).
		Order("post_count DESC, u.id ASC").
		Scan(&rows).Error
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
	err := r.q.WithContext(ctx).
		Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate, Options: clause.LockingOptionsSkipLocked}).
		Where("email_verified = ?", false).
		Order("created_at ASC, id ASC").
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("claim unverified users: %w", err)
	}
	return usersToDomain(rows), nil
}

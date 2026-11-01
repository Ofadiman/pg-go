package repository

import "context"

type Store interface {
	Users() UserRepository
	Posts() PostRepository
	Comments() CommentRepository
	InTx(ctx context.Context, fn func(tx Store) error) error
}

type UserRepository interface {
	Create(ctx context.Context, in NewUser) (User, error)
	CreateMany(ctx context.Context, in []NewUser) ([]User, error)
	GetByID(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, id string, in UserUpdate) (User, error)
	Delete(ctx context.Context, id string) (User, error)
	Upsert(ctx context.Context, in NewUser) (User, error)
	List(ctx context.Context, params ListUsersParams) (Page[User], error)
	ListEligibleAdults(ctx context.Context) ([]User, error)
	UnverifyMinors(ctx context.Context) ([]User, error)
	TopAgeGroups(ctx context.Context, limit int) ([]AgeGroup, error)
	PostStats(ctx context.Context, minPosts int) ([]UserPostStats, error)
	ClaimUnverified(ctx context.Context, limit int) ([]User, error)
}

type PostRepository interface {
	CreateMany(ctx context.Context, in []NewPost) ([]Post, error)
	GetWithAuthor(ctx context.Context, id string) (PostWithAuthor, error)
	ListRecentWithComments(ctx context.Context, limit int) ([]PostWithComments, error)
}

type CommentRepository interface {
	Thread(ctx context.Context, rootID string) ([]ThreadComment, error)
}

package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) (int, error)
	FindOne(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindAllWithConditions(ctx context.Context, conditions string) ([]User, error)
	Update(ctx context.Context, user *User) error
}

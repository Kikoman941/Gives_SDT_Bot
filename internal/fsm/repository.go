package fsm

import "context"

type Repository interface {
	InsertOrUpdate(ctx context.Context, us *UserState) error
	FindOneWithConditions(ctx context.Context, us *UserState, conditions string) error
}

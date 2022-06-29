package fsm

import "context"

type Repository interface {
	InsertOrUpdate(ctx context.Context, us *UserState) error
	FindAllWithConditions(ctx context.Context, conditions string) ([]UserState, error)
}

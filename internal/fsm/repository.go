package fsm

import "context"

type Repository interface {
	InsertOrUpdate(ctx context.Context, us *Userstate) error
	FindAllWithConditions(ctx context.Context, conditions string) ([]Userstate, error)
}

package fsm

import "context"

type Repository interface {
	InsertOrUpdate(ctx context.Context, us *UserState) error
}

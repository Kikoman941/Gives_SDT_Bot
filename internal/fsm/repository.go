package fsm

import "context"

type Repository interface {
	UpdateOrInsert(ctx context.Context, us *UserState) error
}

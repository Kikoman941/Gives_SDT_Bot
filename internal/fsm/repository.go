package fsm

import "context"

type Repository interface {
	Create(ctx context.Context, us UserState) error
	Update(ctx context.Context, us UserState) error
}

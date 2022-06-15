package fsm

import "context"

type Storage interface {
	Create(ctx context.Context, us UserState) error
	Update(ctx context.Context, us UserState) error
}

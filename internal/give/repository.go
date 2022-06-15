package give

import "context"

type Repository interface {
	Create(ctx context.Context, give *Give) (int, error)
	FindOne(ctx context.Context, give *Give) error
	Update(ctx context.Context, give *Give) error
}

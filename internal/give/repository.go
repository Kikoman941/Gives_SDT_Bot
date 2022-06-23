package give

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, give *Give) error
	FindAllWithConditions(ctx context.Context, conditions string) ([]Give, error)
	Update(ctx context.Context, give *Give) error
}

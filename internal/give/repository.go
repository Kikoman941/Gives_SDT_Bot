package give

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, give *Give) error
	FindAllWithConditions(ctx context.Context, conditions string) ([]Give, error)
	UpdateWithConditions(ctx context.Context, conditions string, update string) error
}

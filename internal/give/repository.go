package give

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, give *Give) error
	FindAllWithConditions(ctx context.Context, conditions string, params ...interface{}) ([]Give, error)
	UpdateWithConditions(ctx context.Context, conditions string, update string, params ...interface{}) error
}

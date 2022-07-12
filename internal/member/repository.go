package member

import "context"

type Repository interface {
	Create(ctx context.Context, member *Member) error
	FindRandomLimitWithConditions(ctx context.Context, limit int, conditions string, params ...interface{}) ([]Member, error)
}

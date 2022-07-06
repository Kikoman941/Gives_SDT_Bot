package member

import "context"

type Repository interface {
	Create(ctx context.Context, member *Member) error
}

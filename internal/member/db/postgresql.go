package db

import (
	"Gives_SDT_Bot/internal/member"
	"Gives_SDT_Bot/pkg/client/postgresql"
	"Gives_SDT_Bot/pkg/logging"
	"context"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, member *member.Member) error {
	query := r.client.ModelContext(ctx, member)
	_, err := query.Insert()
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindRandomLimitWithConditions(ctx context.Context, limit int, conditions string, params ...interface{}) ([]member.Member, error) {
	var members []member.Member
	err := r.client.ModelContext(ctx, &members).
		Where(conditions, params...).
		OrderExpr("random()").
		Limit(limit).
		Select()
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (r *repository) DeleteWithConditions(ctx context.Context, conditions string, params ...interface{}) error {
	_, err := r.client.ModelContext(ctx, &member.Member{}).
		Where(conditions, params...).
		Delete()
	if err != nil {
		return err
	}

	return nil
}

func NewRepository(dbClient postgresql.Client, logger *logging.Logger) member.Repository {
	return &repository{
		client: dbClient,
		logger: logger,
	}
}

package db

import (
	"Gives_SDT_Bot/internal/give"
	"Gives_SDT_Bot/pkg/client/postgresql"
	"Gives_SDT_Bot/pkg/logging"
	"context"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, give *give.Give) error {
	query := r.client.ModelContext(ctx, give)
	_, err := query.Insert()
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindOne(ctx context.Context, give *give.Give) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Update(ctx context.Context, give *give.Give) error {
	//TODO implement me
	panic("implement me")
}

func NewRepository(dbClient postgresql.Client, logger *logging.Logger) give.Repository {
	return &repository{
		client: dbClient,
		logger: logger,
	}
}

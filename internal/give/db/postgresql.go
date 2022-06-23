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

func (r *repository) FindAllWithConditions(ctx context.Context, conditions string) ([]give.Give, error) {
	var gives []give.Give

	err := r.client.ModelContext(ctx, &gives).
		Where(conditions).
		Select()
	if err != nil {
		return nil, err
	}

	return gives, nil
}

func (r *repository) UpdateWithConditions(ctx context.Context, conditions string, update string) error {
	var g give.Give
	_, err := r.client.ModelContext(ctx, &g).
		Where(conditions).
		Set(update).
		Update()
	if err != nil {
		return err
	}

	return nil
}

func NewRepository(dbClient postgresql.Client, logger *logging.Logger) give.Repository {
	return &repository{
		client: dbClient,
		logger: logger,
	}
}

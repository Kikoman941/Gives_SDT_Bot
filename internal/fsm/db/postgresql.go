package db

import (
	"Gives_SDT_Bot/internal/fsm"
	"Gives_SDT_Bot/pkg/client/postgresql"
	"Gives_SDT_Bot/pkg/logging"
	"context"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) InsertOrUpdate(ctx context.Context, us *fsm.Userstate) error {
	query := r.client.ModelContext(ctx, us)
	_, err := query.OnConflict("(\"userId\") DO UPDATE").
		Insert()
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAllWithConditions(ctx context.Context, conditions string) ([]fsm.Userstate, error) {
	var states []fsm.Userstate

	err := r.client.ModelContext(ctx, &states).
		Where(conditions).
		Select()
	if err != nil {
		return nil, err
	}
	return states, nil
}

func NewRepository(dbClient postgresql.Client, logger *logging.Logger) fsm.Repository {
	return &repository{
		client: dbClient,
		logger: logger,
	}
}

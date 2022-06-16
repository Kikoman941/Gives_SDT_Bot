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

func (r *repository) UpdateOrInsert(ctx context.Context, us *fsm.UserState) error {
	query := r.client.ModelContext(ctx, us)
	_, err := query.OnConflict("(UserID) DO_UPDATE").
		Set("(State) = State").
		Insert()
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(dbClient postgresql.Client, logger *logging.Logger) fsm.Repository {
	return &repository{
		client: dbClient,
		logger: logger,
	}
}

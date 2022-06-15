package db

import (
	"Gives_SDT_Bot/internal/fsm"
	"Gives_SDT_Bot/pkg/logging"
	"context"
	"github.com/go-pg/pg/v10"
)

type repository struct {
	client *pg.DB
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, us fsm.UserState) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Update(ctx context.Context, us fsm.UserState) error {
	//TODO implement me
	panic("implement me")
}

func NewStorage(dbClient *pg.DB, logger *logging.Logger) fsm.Repository {
	return &repository{
		client: dbClient,
		logger: logger,
	}
}

package db

import (
	"Gives_SDT_Bot/internal/give"
	"Gives_SDT_Bot/pkg/logging"
	"context"
	"github.com/go-pg/pg/v10"
)

type repository struct {
	client *pg.DB
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, give *give.Give) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) FindOne(ctx context.Context, give *give.Give) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Update(ctx context.Context, give *give.Give) error {
	//TODO implement me
	panic("implement me")
}

func NewStorage(dbClient *pg.DB, logger *logging.Logger) give.Repository {
	return &repository{
		client: dbClient,
		logger: logger,
	}
}

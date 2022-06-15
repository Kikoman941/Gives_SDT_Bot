package db

import (
	"Gives_SDT_Bot/internal/give"
	"Gives_SDT_Bot/pkg/logging"
	"context"
	"github.com/go-pg/pg/v10"
)

type db struct {
	db     *pg.DB
	logger *logging.Logger
}

func (d db) Create(ctx context.Context, give *give.Give) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (d db) FindOne(ctx context.Context, give *give.Give) error {
	//TODO implement me
	panic("implement me")
}

func (d db) Update(ctx context.Context, give *give.Give) error {
	//TODO implement me
	panic("implement me")
}

func NewStorage(database *pg.DB, logger *logging.Logger) give.Storage {
	return &db{
		db:     database,
		logger: logger,
	}
}

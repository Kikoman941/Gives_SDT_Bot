package db

import (
	"Gives_SDT_Bot/internal/user"
	"Gives_SDT_Bot/pkg/logging"
	"context"
	"github.com/go-pg/pg/v10"
)

type db struct {
	db     *pg.DB
	logger *logging.Logger
}

func (d *db) Create(ctx context.Context, user *user.User) (int, error) {
	query := d.db.Model(user)
	_, err := query.OnConflict("DO_NOTHING").Insert()
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (d *db) FindOne(ctx context.Context, user *user.User) error {
	query := d.db.Model(user)

	err := query.Select()
	if err != nil {
		return err
	}
	return nil
}

func (d *db) Update(ctx context.Context, user *user.User) error {
	query := d.db.Model(user)

	_, err := query.Update()
	if err != nil {
		return err
	}
	return nil
}

func NewStorage(database *pg.DB, logger *logging.Logger) user.Storage {
	return &db{
		db:     database,
		logger: logger,
	}
}

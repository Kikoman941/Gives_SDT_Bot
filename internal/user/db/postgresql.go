package db

import (
	"Gives_SDT_Bot/internal/user"
	"Gives_SDT_Bot/pkg/logging"
	"context"
	"github.com/go-pg/pg/v10"
)

type repository struct {
	client *pg.DB
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, user *user.User) (int, error) {
	query := r.client.Model(user)
	_, err := query.OnConflict("DO_NOTHING").Insert()
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (r *repository) FindOne(ctx context.Context, user *user.User) error {
	query := r.client.Model(user)

	err := query.Select()
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Update(ctx context.Context, user *user.User) error {
	query := r.client.Model(user)

	_, err := query.Update()
	if err != nil {
		return err
	}
	return nil
}

func NewStorage(dbClient *pg.DB, logger *logging.Logger) user.Repository {
	return &repository{
		client: dbClient,
		logger: logger,
	}
}

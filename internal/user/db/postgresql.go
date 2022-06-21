package db

import (
	"Gives_SDT_Bot/internal/user"
	"Gives_SDT_Bot/pkg/client/postgresql"
	"Gives_SDT_Bot/pkg/logging"
	"context"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, user *user.User) (int, error) {
	query := r.client.ModelContext(ctx, user)
	_, err := query.OnConflict("DO NOTHING").Insert()
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (r *repository) FindAllWithConditions(ctx context.Context, conditions string) ([]user.User, error) {
	var users []user.User

	err := r.client.ModelContext(ctx, &users).
		Where(conditions).
		Select()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) FindOne(ctx context.Context, user *user.User) error {
	query := r.client.ModelContext(ctx, user)

	err := query.Select()
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Update(ctx context.Context, user *user.User) error {
	query := r.client.ModelContext(ctx, user)

	_, err := query.Update()
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(dbClient postgresql.Client, logger *logging.Logger) user.Repository {
	return &repository{
		client: dbClient,
		logger: logger,
	}
}

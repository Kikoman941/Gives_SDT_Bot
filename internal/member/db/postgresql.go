package db

import (
	"Gives_SDT_Bot/internal/member"
	"Gives_SDT_Bot/pkg/client/postgresql"
	"Gives_SDT_Bot/pkg/logging"
	"context"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, member *member.Member) error {
	r.logger.Debug(member)
	query := r.client.ModelContext(ctx, member)
	_, err := query.Insert()
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(dbClient postgresql.Client, logger *logging.Logger) member.Repository {
	return &repository{
		client: dbClient,
		logger: logger,
	}
}

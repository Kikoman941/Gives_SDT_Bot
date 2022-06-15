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

func (r *repository) UpdateOrInsert(ctx context.Context, us fsm.UserState) error {
	query := r.client.Model(us)
	_, err := query.OnConflict("(UserID) DO_UPDATE").
		Set("(State) = State").
		Insert()
	if err != nil {
		return err
	}
	return nil
}

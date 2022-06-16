package postgresql

import (
	"context"
	"errors"
	"github.com/go-pg/pg/v10"
)

type Client interface {
	ModelContext(c context.Context, model ...interface{}) *pg.Query
}

func NewClient(ctx context.Context, dsn string) (*pg.DB, error) {
	opt, err := pg.ParseURL(dsn)
	if err != nil {
		return nil, errors.New("cannot parse postgres dsn")
	}

	db := pg.Connect(opt)

	if err := db.Ping(ctx); err != nil {
		return nil, errors.New("cannot ping postgres db")
	}

	return db, nil
}

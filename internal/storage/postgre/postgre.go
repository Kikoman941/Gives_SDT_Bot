package postgre

import (
	"fmt"
	"github.com/go-pg/pg/v10"
)

type PostgreDB struct {
	db *pg.DB
}

func NewPostgresDB(dsn string) (*PostgreDB, error) {
	opt, err := pg.ParseURL(dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot parse storage dsn: %s", err)
	}

	db := pg.Connect(opt)

	return &PostgreDB{
		db: db,
	}, nil
}

func (pg *PostgreDB) Select(model interface{}, condition string) error {
	query := pg.db.Model(model)

	if condition != "" {
		query.Where(condition)
	}

	err := query.Select()
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgreDB) Insert(model interface{}) (interface{}, error) {
	_, err := pg.db.Model(model).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return nil, err
	}
	return model, nil
}

package postgre

import (
	"fmt"
	"github.com/go-pg/pg/v10"
)

type PostgresDB struct {
	db *pg.DB
}

func NewPostgresDB(dsn string) (*PostgresDB, error) {
	opt, err := pg.ParseURL(dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot parse storage dsn: %s", err)
	}

	db := pg.Connect(opt)

	return &PostgresDB{
		db: db,
	}, nil
}

func (pg *PostgresDB) Select(model interface{}, tableName string, condition string) error {
	query := pg.db.Model(model)

	if tableName != "" {
		query.Table(tableName)
	}

	if condition != "" {
		query.Where(condition)
	}

	err := query.Select()
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresDB) Insert(model interface{}) (interface{}, error) {
	_, err := pg.db.Model(model).SelectOrInsert()
	if err != nil {
		return nil, err
	}
	return model, nil
}

package db

import (
	"database/sql"
	"errors"
	"fmt"

	c "github.com/boginskiy/Clicki/cmd/config"
	l "github.com/boginskiy/Clicki/internal/logger"
)

type StoreDB struct {
	Logger l.Logger
	DB     *sql.DB
}

func NewStoreDB(kwargs c.VarGetter, logger l.Logger) (*StoreDB, error) {
	db, err := sql.Open("postgres", kwargs.GetDB())

	// TODO!
	l := fmt.Sprintf("PingDB:%v\n", db.Ping())
	logger.RaiseError(errors.New("POSTGRES3"), l, nil)

	if err != nil {
		return nil, err
	}
	return &StoreDB{
		Logger: logger,
		DB:     db,
	}, nil
}

func (d *StoreDB) CloseDB() {
	d.DB.Close()
}

func (d *StoreDB) GetDB() *sql.DB {
	return d.DB
}

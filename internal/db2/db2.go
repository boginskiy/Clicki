package db2

import (
	"database/sql"

	c "github.com/boginskiy/Clicki/cmd/config"
	l "github.com/boginskiy/Clicki/internal/logger"
)

type ConnDB struct {
	Logger l.Logger
	DB     *sql.DB
}

// NewDBConn - Развертывание PostgreSQL
func NewConnDB(kwargs c.VarGetter, logger l.Logger) *ConnDB {
	db, err := sql.Open("postgres", kwargs.GetDB())
	if err != nil {
		logger.RaiseFatal(err, "NewDBConn>sql.Open", l.Fields{"data": kwargs.GetDB()})
	}
	return &ConnDB{
		Logger: logger,
		DB:     db,
	}
}

func (d *ConnDB) CloseDB() {
	d.DB.Close()
}

func (d *ConnDB) GetDB() *sql.DB {
	return d.DB
}

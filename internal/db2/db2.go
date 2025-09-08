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
	DATABASE_DSN := "postgres://videos:userpassword@localhost:5432/videos?sslmode=disable"

	db, err := sql.Open("postgres", DATABASE_DSN)
	if err != nil {
		logger.RaiseFatal("NewDBConn>sql.Open", l.Fields{"data": DATABASE_DSN})
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

// Отдельные кварги на подключение DB
// DATABASE_DSN
// -d

package db

import (
	"database/sql"

	c "github.com/boginskiy/Clicki/cmd/config"
	l "github.com/boginskiy/Clicki/internal/logger"
	_ "github.com/lib/pq"
)

type StoreDB struct {
	Logger l.Logger
	DB     *sql.DB
}

func NewStoreDB(kwargs c.VarGetter, logger l.Logger) (*StoreDB, error) {
	db, err := sql.Open("postgres", kwargs.GetDB())
	if err != nil {
		return nil, err
	}

	// "Самостоятельное" создание таблиц"
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			original_url TEXT NOT NULL,
			short_url VARCHAR(8) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`)
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

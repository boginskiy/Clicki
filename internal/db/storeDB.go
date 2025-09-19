package db

import (
	"database/sql"

	c "github.com/boginskiy/Clicki/cmd/config"
	l "github.com/boginskiy/Clicki/internal/logger"
	_ "github.com/lib/pq"
)

func createUrls(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS urls (
						id SERIAL PRIMARY KEY,
						correlation_id TEXT,
						original_url TEXT NOT NULL UNIQUE,
						short_url TEXT NOT NULL,
						created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`)
	return err
}

type StoreDB struct {
	Logger l.Logger
	DB     *sql.DB
}

func NewStoreDB(kwargs c.VarGetter, logger l.Logger) (*StoreDB, error) {
	db, err := sql.Open("postgres", kwargs.GetDB())
	if err != nil {
		return nil, err
	}

	err = createUrls(db)
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

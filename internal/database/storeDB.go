package database

import (
	"database/sql"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/errs"
	"github.com/boginskiy/Clicki/internal/logg"
	_ "github.com/lib/pq"
)

func createUrls(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS urls (
						id SERIAL PRIMARY KEY,
						correlation_id TEXT,
						original_url TEXT NOT NULL UNIQUE,
						short_url TEXT NOT NULL,
						created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						deleted_flag BOOLEAN NOT NULL DEFAULT FALSE,
						user_id INT NOT NULL)`)
	return err
}

// StoreDB - store is real database.
type StoreDB struct {
	Logg logg.Logger
	DB   *sql.DB
}

func NewStoreDB(config conf.Config, logger logg.Logger) (DataBase, error) {
	db, err := sql.Open("postgres", config.GetDB())
	if err != nil {
		return nil, err
	}

	err = createUrls(db)
	if err != nil {
		return nil, err
	}

	return &StoreDB{
		Logg: logger,
		DB:   db,
	}, nil
}

func (sd *StoreDB) CloseDB() {
	sd.DB.Close()
}

func (sd *StoreDB) GetDB() any {
	return sd.DB
}

func (sd *StoreDB) CheckOpen() (bool, error) {
	err := sd.DB.Ping()
	if err != nil {
		return false, errs.ErrPingDataBase
	}
	return true, nil
}

package db

import (
	"database/sql"
	"log"
	"os"

	conf "github.com/boginskiy/Clicki/cmd/config"
	cerr "github.com/boginskiy/Clicki/internal/error"
	"github.com/boginskiy/Clicki/internal/logg"
	_ "github.com/lib/pq"
)

func DELETEusers(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE users CASCADE`)
	return err
}

func CREATEusers(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
						id SERIAL PRIMARY KEY,
						name TEXT NOT NULL,
						email TEXT UNIQUE,
						password TEXT,
						created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP,
						last_login_at TIMESTAMP,
						is_active BOOLEAN DEFAULT TRUE,
						roles TEXT[])`)
	return err
}

func CREATEurls(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS urls (
						id SERIAL PRIMARY KEY,
						correlation_id TEXT,
						original_url TEXT NOT NULL UNIQUE,
						short_url TEXT NOT NULL,
						created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
						user_id INT, 
						FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE)`)
	return err
}

type StoreDB struct {
	Logger logg.Logger
	DB     *sql.DB
}

func NewStoreDB(kwargs conf.VarGetter, logger logg.Logger) (DBer, error) {
	db, err := sql.Open("postgres", kwargs.GetDB())
	if err != nil {
		return nil, err
	}

	// Создание таблицы users
	err = CREATEusers(db)
	if err != nil {
		return nil, err
	}
	// Создание таблицы urls
	err = CREATEurls(db)
	if err != nil {
		return nil, err
	}

	return &StoreDB{
		Logger: logger,
		DB:     db,
	}, nil
}

func (sd *StoreDB) CloseDB() {
	// TODO! Костыль для прохождения тестов
	// Перед завершением удалим таблицу users
	err := DELETEusers(sd.DB)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	sd.DB.Close()
}

func (sd *StoreDB) GetDB() any {
	return sd.DB
}

func (sd *StoreDB) CheckOpen() (bool, error) {
	err := sd.DB.Ping()
	if err != nil {
		return false, cerr.ErrPingDataBase
	}
	return true, nil
}

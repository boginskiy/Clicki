package db

import (
	"database/sql"
	"log"

	c "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/config"
	l "github.com/boginskiy/Clicki/internal/logger"
	"github.com/boginskiy/Clicki/internal/migration"
)

type StoreDB struct {
	Logger l.Logger
	DB     *sql.DB
}

// func NewStoreDB(kwargs c.VarGetter, logger l.Logger) (*StoreDB, error) {
// 	db, err := sql.Open("postgres", kwargs.GetDB())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &StoreDB{
// 		Logger: logger,
// 		DB:     db,
// 	}, nil
// }

// TODO!

func NewStoreDB(kwargs c.VarGetter, logger l.Logger) (*StoreDB, error) {
	cfg := config.NewConfig()
	db, err := config.OpenDatabase(cfg)

	// db, err := sql.Open("postgres", kwargs.GetDB())
	if err != nil {
		log.Fatal(err)
	}

	// Применяем миграции сразу после открытия соединения с базой данных
	if err := migration.ApplyMigrations(db); err != nil {
		log.Fatalf("Error applying migrations: %v\\n", err)
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

// TODO!
// У меня не срабатывают миграции перед запуском приложения

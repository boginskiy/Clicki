package main

import (
	"github.com/boginskiy/Clicki/cmd/config"

	"github.com/boginskiy/Clicki/internal/db"

	"github.com/boginskiy/Clicki/internal/logg"
	repo "github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/internal/server"
)

// ChoiceOfCreatingDB - Функции для создания слоя DB
func ChoiceOfCreatingDB(kwargs config.VarGetter) func(config.VarGetter, logg.Logger) (db.DBer, error) {
	if kwargs.GetDB() != "" {
		return db.NewStoreDB
	} else if kwargs.GetPathToStore() != "" {
		return db.NewStoreFile
	} else {
		return db.NewStoreMap
	}
}

// ChoiceOfRepositoryDB - Функции для создания слоя Repo
func ChoiceOfRepositoryDB(database db.DBer) func(config.VarGetter, db.DBer) (repo.Repository, error) {
	switch database.(type) {
	case *db.StoreDB:
		return repo.NewRepositoryDBURL
	case *db.StoreFile:
		return repo.NewRepositoryFileURL
	case *db.StoreMap:
		return repo.NewRepositoryMapURL
	default:
		return nil
	}
}

func main() {
	// Logger
	BaseLogger := logg.NewLogg("LogBase.log", "FATAL")

	// Kwargs
	Variables := config.NewVariables(BaseLogger)

	// Database
	newDataBase := ChoiceOfCreatingDB(Variables)
	database, err := newDataBase(Variables, BaseLogger)
	if err != nil {
		BaseLogger.RaiseFatal(err, "main>NewDataBase", nil)
	}

	// Repository
	NewRepository := ChoiceOfRepositoryDB(database)
	Repository, err := NewRepository(Variables, database)
	if err != nil {
		BaseLogger.RaiseFatal(err, "main>NewRepository", nil)
	}

	defer BaseLogger.CloseDesc()
	defer database.CloseDB()

	server.Run(Variables, BaseLogger, Repository)
}

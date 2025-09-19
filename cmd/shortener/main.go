package main

import (
	"github.com/boginskiy/Clicki/cmd/config"

	"github.com/boginskiy/Clicki/internal/db"

	l "github.com/boginskiy/Clicki/internal/logger"
	rp "github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/internal/server"
)

func SetDB(kwargs config.VarGetter, logger l.Logger) db.DBer {
	if kwargs.GetDB() != "" {
		// Используем для хранения данных DataBase
		db, err := db.NewStoreDB(kwargs, logger)
		logger.RaiseFatal(err, "SetDB>NewStoreDB", nil)
		return db

	} else if kwargs.GetPathToStore() != "" {
		// Используем для хранения данных файл
		db, err := db.NewStoreFile(kwargs, logger)
		logger.RaiseFatal(err, "SetDB>NewStoreFile", nil)
		return db

	} else {
		// Используем для хранения данных Map
		db, err := db.NewStoreMap(kwargs, logger)
		logger.RaiseFatal(err, "SetDB>NewStoreMap", nil)
		return db
	}
}

func main() {
	BaseLogger := l.NewLogg("LogBase.log", "FATAL")
	Variables := config.NewVariables(BaseLogger)
	Database := SetDB(Variables, BaseLogger)

	defer BaseLogger.CloseDesc()
	defer Database.CloseDB()

	// Заворачиваем в слой Repository DataBase
	Repo, ok := Database.(rp.URLRepository)
	if !ok {
		Repo = rp.NewSQLURLRepository(Variables, Database)
	}

	server.Run(Variables, BaseLogger, Repo)
}

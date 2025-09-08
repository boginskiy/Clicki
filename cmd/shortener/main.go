package main

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db2"
	"github.com/boginskiy/Clicki/internal/logger"
	"github.com/boginskiy/Clicki/internal/server"
)

func main() {
	BaseLogger := logger.NewLogg("LogBase.log", "FATAL")
	Variables := config.NewVariables(BaseLogger)
	Database := db2.NewConnDB(Variables, BaseLogger)

	defer BaseLogger.CloseDesc()
	defer Database.CloseDB()

	server.Run(Variables, BaseLogger, Database)
}

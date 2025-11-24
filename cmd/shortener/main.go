package main

import (
	"github.com/boginskiy/Clicki/cmd/config"

	"github.com/boginskiy/Clicki/cmd/server"
	"github.com/boginskiy/Clicki/internal/logg"
)

func main() {
	logger := logg.NewLogg("base.log", "FATAL")
	config := config.NewVariables(logger)
	layers := server.NewLayers(config, logger)

	layers.NewLayerDB()
	repo := layers.NewLayerRepo()

	defer logger.Close()
	defer layers.Close()

	server.Run(config, logger, repo)
}

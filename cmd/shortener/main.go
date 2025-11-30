package main

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/app"
	"github.com/boginskiy/Clicki/internal/logg"
)

// main - start of Appl.
func main() {
	logger := logg.NewLogg("main.log", "FATAL")
	cfg := config.NewVariables(logger)

	app.NewApp(cfg, logger).Start()
}

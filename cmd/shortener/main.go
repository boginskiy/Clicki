package main

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/app"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/pkg"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// main - start of Appl.
func main() {
	pkg.PrintInfo(buildVersion, buildDate, buildCommit)

	logger := logg.NewLogg("main.log", "FATAL")
	cfg := config.NewVariables(logger)

	app.NewApp(cfg, logger).Start()
}

package testinit

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
)

func InitConfAndLogg(pathToLogg string) (*config.Variables, *logg.Logg) {
	// Init.
	logg := logg.NewLogg(pathToLogg, "ERROR")
	config := InitConfig()
	return config, logg
}

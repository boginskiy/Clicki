package tserv

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
)

func InitConfig() *config.Variables {
	return &config.Variables{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080",
		ArgsCLI:       &config.ArgsCLI{},
		ArgsENV: &config.ArgsENV{
			SoftDeleteTime: 10,
			HardDeleteTime: 20,
		},
	}
}

func InitConfAndLogg(pathToLogg string) (*config.Variables, *logg.Logg) {
	// Init.
	logg := logg.NewLogg(pathToLogg, "INFO")
	config := InitConfig()
	return config, logg
}

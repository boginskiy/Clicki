package config

import (
	"log"

	"github.com/caarlos0/env"
)

type ArgsEnviron struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	PathToStore   string `env:"FILE_STORAGE_PATH"`
	DB            string `env:"DATABASE_DSN"`
	BaseURL       string `env:"BASE_URL"`
	LogFile       string `env:"LOG_FILE"`
}

func NewArgsEnviron() *ArgsEnviron {
	ArgsEnv := new(ArgsEnviron)
	ArgsEnv.ParseFlags()
	return ArgsEnv
}

func (e *ArgsEnviron) ParseFlags() {
	err := env.Parse(e)
	if err != nil {
		log.Fatal(err)
	}
	// Default
	if e.LogFile == "" {
		e.LogFile = "LogInfo.log"
	}

}

func (e *ArgsEnviron) GetSrvAddr() (ServerAddress string) {
	return e.ServerAddress
}

func (e *ArgsEnviron) GetBaseURL() (BaseURL string) {
	return e.BaseURL
}

func (e *ArgsEnviron) GetPathToStore() (PathToStore string) {
	return e.PathToStore
}

func (e *ArgsEnviron) GetDB() (DB string) {
	return e.DB
}

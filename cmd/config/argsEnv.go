package config

import (
	"log"

	"github.com/caarlos0/env"
)

type ArgsEnviron struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
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
}

func (e *ArgsEnviron) GetSrvAddr() (ServerAddress string) {
	return e.ServerAddress
}

func (e *ArgsEnviron) GetBaseUrl() (BaseURL string) {
	return e.BaseURL
}

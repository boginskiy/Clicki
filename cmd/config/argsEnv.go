package config

import (
	"log"

	"github.com/caarlos0/env"
)

type ArgsEnviron struct {
	Server_address string `env:"SERVER_ADDRESS"`
	Base_url       string `env:"BASE_URL"`
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

func (e *ArgsEnviron) GetSrvAddr() (Server_address string) {
	return e.Server_address
}

func (e *ArgsEnviron) GetBaseUrl() (Base_url string) {
	return e.Base_url
}

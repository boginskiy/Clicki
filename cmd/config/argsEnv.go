package config

import (
	"log"

	"github.com/caarlos0/env"
)

type ArgsEnviron struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	NameLogInfo   string `env:"NAME_LOG_INFO"`
	NameLogFatal  string `env:"NAME_LOG_FATAL"`
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
	if e.NameLogInfo == "" {
		e.NameLogInfo = "LogInfo.log"
	}
	if e.NameLogFatal == "" {
		e.NameLogFatal = "LogFatal.log"
	}
}

func (e *ArgsEnviron) GetSrvAddr() (ServerAddress string) {
	return e.ServerAddress
}

func (e *ArgsEnviron) GetBaseURL() (BaseURL string) {
	return e.BaseURL
}

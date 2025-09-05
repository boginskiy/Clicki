package config

import (
	"log"
	"os"
	"strings"

	"github.com/caarlos0/env"
)

type ArgsENV struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	PathToStore   string `env:"FILE_STORAGE_PATH"`
	BaseURL       string `env:"BASE_URL"`
	NameLogInfo   string `env:"NAME_LOG_INFO"`
	NameLogFatal  string `env:"NAME_LOG_FATAL"`
}

func NewArgsENV() *ArgsENV {
	args := new(ArgsENV)
	args.ParseFlags()
	return args
}

func (e *ArgsENV) ParseFlags() {
	err := env.Parse(e)
	if err != nil {
		log.Fatal(err)
	}
	// Default
	valueLogInfo := strings.TrimSpace(os.Getenv("NAME_LOG_INFO"))
	valueLogFatal := strings.TrimSpace(os.Getenv("NAME_LOG_FATAL"))

	if len(valueLogInfo) == 0 {
		e.NameLogInfo = "LogInfo.log"
	}
	if len(valueLogFatal) == 0 {
		e.NameLogFatal = "LogFatal.log"
	}
}

func (e *ArgsENV) GetSrvAddr() (ServerAddress string) {
	return e.ServerAddress
}

func (e *ArgsENV) GetBaseURL() (BaseURL string) {
	return e.BaseURL
}

func (e *ArgsENV) GetPathToStore() (PathToStore string) {
	return e.PathToStore
}

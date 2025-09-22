package config

import (
	"log"
	"os"
	"strings"

	"github.com/caarlos0/env"
)

type ArgsENV struct {
	ServerAddress string `env:"SERVER_ADDRESS"`    //
	PathToStore   string `env:"FILE_STORAGE_PATH"` //
	DB            string `env:"DATABASE_DSN"`      //
	BaseURL       string `env:"BASE_URL"`          //
	LogFile       string `env:"LOG_FILE"`
	MaxRetries    int    `env:"MAX_RETRIES"`
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
	valueLogFile := strings.TrimSpace(os.Getenv("LOG_FILE"))
	if len(valueLogFile) == 0 {
		e.LogFile = "LogInfo.log"
	}

	valueMaxRetries := strings.TrimSpace(os.Getenv("MAX_RETRIES"))
	if len(valueMaxRetries) == 0 {
		e.MaxRetries = 3
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

func (e *ArgsENV) GetDB() (DB string) {
	return e.DB
}

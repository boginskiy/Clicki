package config

import (
	"log"
	"os"
	"strings"

	"github.com/caarlos0/env"
)

type ArgsENV struct {
	ServerAddress  string `env:"SERVER_ADDRESS"`    //
	PathToStore    string `env:"FILE_STORAGE_PATH"` //
	DB             string `env:"DATABASE_DSN"`      //
	BaseURL        string `env:"BASE_URL"`          //
	LogFile        string `env:"LOG_FILE"`
	MaxRetries     int    `env:"MAX_RETRIES"`
	TokenLiveTime  int    `env:"TOKEN_LIVE_TIME"`
	CokiLiveTime   int    `env:"COKI_LIVE_TIME"`
	NameCoki       string `env:"NAME_COKI"`
	SecretKey      string `env:"SECRET_KEY"`
	SoftDeleteTime int    `env:"SOFT_DELETE_TIME"`
	HardDeleteTime int    `env:"HARD_DELETE_TIME"`
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
	valueStr := strings.TrimSpace(os.Getenv("LOG_FILE"))
	if len(valueStr) == 0 {
		e.LogFile = "LogInfo.log"
	}

	valueStr = strings.TrimSpace(os.Getenv("MAX_RETRIES"))
	if len(valueStr) == 0 {
		e.MaxRetries = 3
	}

	valueStr = strings.TrimSpace(os.Getenv("TOKEN_LIVE_TIME"))
	if len(valueStr) == 0 {
		e.TokenLiveTime = 10
	}

	valueStr = strings.TrimSpace(os.Getenv("COKI_LIVE_TIME"))
	if len(valueStr) == 0 {
		e.CokiLiveTime = 300
	}

	valueStr = strings.TrimSpace(os.Getenv("NAME_COKI"))
	if len(valueStr) == 0 {
		e.NameCoki = "auth_token"
	}

	valueStr = strings.TrimSpace(os.Getenv("SECRET_KEY"))
	if len(valueStr) == 0 {
		e.SecretKey = "Ld5pS4Gw"
	}

	valueStr = strings.TrimSpace(os.Getenv("SOFT_DELETE_TIME"))
	if len(valueStr) == 0 {
		e.SoftDeleteTime = 10
	}

	valueStr = strings.TrimSpace(os.Getenv("HARD_DELETE_TIME"))
	if len(valueStr) == 0 {
		e.HardDeleteTime = 60
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

package config

import (
	"log"
	"os"
	"strings"

	"github.com/caarlos0/env"
)

// ArgsENV - struct for args environment.
// generate:reset
type ArgsENV struct {
	ServerAddress  string `env:"SERVER_ADDRESS"`
	PathToStore    string `env:"FILE_STORAGE_PATH"`
	DB             string `env:"DATABASE_DSN"`
	BaseURL        string `env:"BASE_URL"`
	AuditFile      string `env:"AUDIT_FILE"`
	AuditURL       string `env:"AUDIT_URL"`
	LogFile        string `env:"LOG_FILE"`
	MaxRetries     int    `env:"MAX_RETRIES"`
	TokenLiveTime  int    `env:"TOKEN_LIVE_TIME"`
	CokiLiveTime   int    `env:"COKI_LIVE_TIME"`
	NameCoki       string `env:"NAME_COKI"`
	SecretKey      string `env:"SECRET_KEY"`
	SoftDeleteTime int    `env:"SOFT_DELETE_TIME"`
	HardDeleteTime int    `env:"HARD_DELETE_TIME"`
	EnableHTTPS    string `env:"ENABLE_HTTPS"`
}

func NewArgsENV() *ArgsENV {
	args := new(ArgsENV)
	args.ParseFlags()
	return args
}

func (ae *ArgsENV) ParseFlags() {
	err := env.Parse(ae)
	if err != nil {
		log.Fatal(err)
	}

	// Default value.
	valueStr := strings.TrimSpace(os.Getenv("LOG_FILE"))
	if len(valueStr) == 0 {
		ae.LogFile = "infra.log"
	}

	valueStr = strings.TrimSpace(os.Getenv("MAX_RETRIES"))
	if len(valueStr) == 0 {
		ae.MaxRetries = 3
	}

	valueStr = strings.TrimSpace(os.Getenv("TOKEN_LIVE_TIME"))
	if len(valueStr) == 0 {
		ae.TokenLiveTime = 10
	}

	valueStr = strings.TrimSpace(os.Getenv("COKI_LIVE_TIME"))
	if len(valueStr) == 0 {
		ae.CokiLiveTime = 300
	}

	valueStr = strings.TrimSpace(os.Getenv("NAME_COKI"))
	if len(valueStr) == 0 {
		ae.NameCoki = "auth_token"
	}

	valueStr = strings.TrimSpace(os.Getenv("SECRET_KEY"))
	if len(valueStr) == 0 {
		ae.SecretKey = "Ld5pS4Gw"
	}

	valueStr = strings.TrimSpace(os.Getenv("SOFT_DELETE_TIME"))
	if len(valueStr) == 0 {
		ae.SoftDeleteTime = 10
	}

	valueStr = strings.TrimSpace(os.Getenv("HARD_DELETE_TIME"))
	if len(valueStr) == 0 {
		ae.HardDeleteTime = 50
	}

}

func (ae *ArgsENV) GetSrvAddr() (ServerAddress string) {
	return ae.ServerAddress
}

func (ae *ArgsENV) GetBaseURL() (BaseURL string) {
	return ae.BaseURL
}

func (ae *ArgsENV) GetPathToStore() (PathToStore string) {
	return ae.PathToStore
}

func (ae *ArgsENV) GetDB() (DB string) {
	return ae.DB
}

func (ae *ArgsENV) GetAuditFile() (AuditFile string) {
	return ae.AuditFile
}

func (ae *ArgsENV) GetAuditURL() (AuditURL string) {
	return ae.AuditURL
}

func (ae *ArgsENV) GetEnableHTTPS() (EnableHTTPS string) {
	return ae.EnableHTTPS
}

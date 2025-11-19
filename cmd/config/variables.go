package config

import (
	"strings"

	"github.com/boginskiy/Clicki/internal/logg"
)

type Variables struct {
	Logger        logg.Logger
	ServerAddress string
	PathToStore   string
	AuditFile     string
	AuditURL      string
	BaseURL       string
	DB            string
	ArgsCLI       *ArgsCLI
	ArgsENV       *ArgsENV
}

func NewVariables(logger logg.Logger) *Variables {
	tmpVar := &Variables{
		Logger:  logger,
		ArgsCLI: NewArgsCLI(),
		ArgsENV: NewArgsENV(),
	}
	tmpVar.extSettingsArgs()
	return tmpVar
}

func (v *Variables) argsTrim(arg string) string {
	return strings.TrimSpace(arg)
}

func (v *Variables) argsPrioryty(envFunc, cliFunc func() string) string {
	arg := v.argsTrim(envFunc())  // Clean arg
	arg2 := v.argsTrim(cliFunc()) // Clean arg

	if len(arg) > 0 {
		return arg
	} else {
		return arg2
	}
}

func (v *Variables) extSettingsArgs() {
	v.PathToStore = v.argsPrioryty(v.ArgsENV.GetPathToStore, v.ArgsCLI.GetPathToStore)
	v.ServerAddress = v.argsPrioryty(v.ArgsENV.GetSrvAddr, v.ArgsCLI.GetSrvAddr)
	v.AuditFile = v.argsPrioryty(v.ArgsENV.GetAuditFile, v.ArgsCLI.GetAuditFile)
	v.AuditURL = v.argsPrioryty(v.ArgsENV.GetAuditURL, v.ArgsCLI.GetAuditURL)
	v.BaseURL = v.argsPrioryty(v.ArgsENV.GetBaseURL, v.ArgsCLI.GetBaseURL)
	v.DB = v.argsPrioryty(v.ArgsENV.GetDB, v.ArgsCLI.GetDB)
}

func (v *Variables) GetSrvAddr() (ServerAddress string) {
	return v.ServerAddress
}

func (v *Variables) GetBaseURL() (BaseURL string) {
	return v.BaseURL
}

func (v *Variables) GetPathToStore() (PathToStore string) {
	return v.PathToStore
}

func (v *Variables) GetDB() (DB string) {
	return v.DB
}

func (v *Variables) GetLogFile() (LogFile string) {
	return v.ArgsENV.LogFile
}

func (v *Variables) GetMaxRetries() (MaxRetries int) {
	return v.ArgsENV.MaxRetries
}

func (v *Variables) GetTokenLiveTime() (TokenLiveTime int) {
	return v.ArgsENV.TokenLiveTime
}

func (v *Variables) GetCokiLiveTime() (CokiLiveTime int) {
	return v.ArgsENV.CokiLiveTime
}

func (v *Variables) GetNameCoki() (NameCoki string) {
	return v.ArgsENV.NameCoki
}

func (v *Variables) GetSecretKey() (SecretKey string) {
	return v.ArgsENV.SecretKey
}

func (v *Variables) GetSoftDeleteTime() (SoftDeleteTime int) {
	return v.ArgsENV.SoftDeleteTime
}

func (v *Variables) GetHardDeleteTime() (HardDeleteTime int) {
	return v.ArgsENV.HardDeleteTime
}

func (v *Variables) GetAuditFile() (AuditFile string) {
	return v.AuditFile
}

func (v *Variables) GetAuditURL() (AuditURL string) {
	return v.AuditURL
}

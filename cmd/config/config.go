package config

import (
	"strings"

	"github.com/boginskiy/Clicki/internal/logg"
)

// Conf - struct about config interface.
type Conf struct {
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

func NewVariables(logger logg.Logger) *Conf {
	tmpVar := &Conf{
		Logger:  logger,
		ArgsCLI: NewArgsCLI(),
		ArgsENV: NewArgsENV(),
	}
	tmpVar.extSettingsArgs()
	return tmpVar
}

func (v *Conf) argsTrim(arg string) string {
	return strings.TrimSpace(arg)
}

func (v *Conf) argsPrioryty(envFunc, cliFunc func() string) string {
	arg := v.argsTrim(envFunc())  // Clean arg
	arg2 := v.argsTrim(cliFunc()) // Clean arg

	if len(arg) > 0 {
		return arg
	} else {
		return arg2
	}
}

func (v *Conf) extSettingsArgs() {
	v.PathToStore = v.argsPrioryty(v.ArgsENV.GetPathToStore, v.ArgsCLI.GetPathToStore)
	v.ServerAddress = v.argsPrioryty(v.ArgsENV.GetSrvAddr, v.ArgsCLI.GetSrvAddr)
	v.AuditFile = v.argsPrioryty(v.ArgsENV.GetAuditFile, v.ArgsCLI.GetAuditFile)
	v.AuditURL = v.argsPrioryty(v.ArgsENV.GetAuditURL, v.ArgsCLI.GetAuditURL)
	v.BaseURL = v.argsPrioryty(v.ArgsENV.GetBaseURL, v.ArgsCLI.GetBaseURL)
	v.DB = v.argsPrioryty(v.ArgsENV.GetDB, v.ArgsCLI.GetDB)
}

func (v *Conf) GetSrvAddr() (ServerAddress string) {
	return v.ServerAddress
}

func (v *Conf) GetBaseURL() (BaseURL string) {
	return v.BaseURL
}

func (v *Conf) GetPathToStore() (PathToStore string) {
	return v.PathToStore
}

func (v *Conf) GetDB() (DB string) {
	return v.DB
}

func (v *Conf) GetLogFile() (LogFile string) {
	return v.ArgsENV.LogFile
}

func (v *Conf) GetMaxRetries() (MaxRetries int) {
	return v.ArgsENV.MaxRetries
}

func (v *Conf) GetTokenLiveTime() (TokenLiveTime int) {
	return v.ArgsENV.TokenLiveTime
}

func (v *Conf) GetCokiLiveTime() (CokiLiveTime int) {
	return v.ArgsENV.CokiLiveTime
}

func (v *Conf) GetNameCoki() (NameCoki string) {
	return v.ArgsENV.NameCoki
}

func (v *Conf) GetSecretKey() (SecretKey string) {
	return v.ArgsENV.SecretKey
}

func (v *Conf) GetSoftDeleteTime() (SoftDeleteTime int) {
	return v.ArgsENV.SoftDeleteTime
}

func (v *Conf) GetHardDeleteTime() (HardDeleteTime int) {
	return v.ArgsENV.HardDeleteTime
}

func (v *Conf) GetAuditFile() (AuditFile string) {
	return v.AuditFile
}

func (v *Conf) GetAuditURL() (AuditURL string) {
	return v.AuditURL
}

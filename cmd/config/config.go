package config

import (
	"strings"

	"github.com/boginskiy/Clicki/internal/logg"
)

// Variables - struct about config interface.
type Variables struct {
	Logger        logg.Logger
	ServerAddress string
	EnableHTTPS   string
	PathToStore   string
	ConfigFile    string
	AuditFile     string
	AuditURL      string
	BaseURL       string
	DB            string
	ArgsCLI       *ArgsCLI
	ArgsENV       *ArgsENV
	ArgsJSON      *ArgsJSON
}

func NewVariables(logger logg.Logger) *Variables {
	vr := &Variables{
		Logger:  logger,
		ArgsCLI: NewArgsCLI(),
		ArgsENV: NewArgsENV(),
	}
	// ArgsJSON
	configFile := vr.argsPrioryty(
		vr.ArgsENV.GetConfigFile, vr.ArgsCLI.GetConfigFile)
	vr.ArgsJSON = NewArgsJSON(configFile)

	vr.settingsPrioryty(vr.ArgsENV, vr.ArgsCLI)
	return vr
}

func (v *Variables) settingsPrioryty(obj1, obj2 ConfigPrioryty) {
	v.PathToStore = v.argsPrioryty(obj1.GetPathToStore, obj2.GetPathToStore)
	v.EnableHTTPS = v.argsPrioryty(obj1.GetEnableHTTPS, obj2.GetEnableHTTPS)
	v.ServerAddress = v.argsPrioryty(obj1.GetSrvAddr, obj2.GetSrvAddr)
	v.AuditFile = v.argsPrioryty(obj1.GetAuditFile, obj2.GetAuditFile)
	v.AuditURL = v.argsPrioryty(obj1.GetAuditURL, obj2.GetAuditURL)
	v.BaseURL = v.argsPrioryty(obj1.GetBaseURL, obj2.GetBaseURL)
	v.DB = v.argsPrioryty(obj1.GetDB, obj2.GetDB)

	if _, ok := obj2.(*ArgsJSON); ok {
		return
	}

	if v.ArgsJSON != nil {
		v.settingsPrioryty(v, v.ArgsJSON)
	}
}

func (v *Variables) argsPrioryty(func1, func2 func() string) string {
	arg := strings.TrimSpace(func1())  // Clean arg
	arg2 := strings.TrimSpace(func2()) // Clean arg

	if len(arg) > 0 {
		return arg
	} else {
		return arg2
	}
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

func (v *Variables) GetEnableHTTPS() (EnableHTTPS string) {
	return v.EnableHTTPS
}

package config

import (
	"strings"
)

type Variables struct {
	ServerAddress string
	PathToStore   string
	BaseURL       string
	ArgsCLI       *ArgsCLI
	ArgsENV       *ArgsENV
}

func NewVariables() *Variables {
	tmpVar := &Variables{
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
	arg := envFunc()      // Get arg
	arg = v.argsTrim(arg) // Clean arg

	if len(arg) > 0 {
		return arg
	} else {
		return cliFunc()
	}
}

func (v *Variables) extSettingsArgs() {
	v.PathToStore = v.argsPrioryty(v.ArgsENV.GetPathToStore, v.ArgsCLI.GetPathToStore)
	v.ServerAddress = v.argsPrioryty(v.ArgsENV.GetSrvAddr, v.ArgsCLI.GetSrvAddr)
	v.BaseURL = v.argsPrioryty(v.ArgsENV.GetBaseURL, v.ArgsCLI.GetBaseURL)
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

func (v *Variables) GetNameLogInfo() (NameLogInfo string) {
	return v.ArgsENV.NameLogInfo
}

func (v *Variables) GetNameLogFatal() (NameLogFatal string) {
	return v.ArgsENV.NameLogFatal
}

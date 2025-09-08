package config

import (
	l "github.com/boginskiy/Clicki/internal/logger"
)

type Variables struct {
	Logger          l.Logger
	ServerAddress   string
	PathToStore     string
	BaseURL         string
	DB              string
	ArgsCommandLine *ArgsCommandLine
	ArgsEnviron     *ArgsEnviron
}

func NewVariables(logger l.Logger) *Variables {
	tmpVar := &Variables{
		Logger:          logger,
		ArgsCommandLine: NewArgsCommandLine(),
		ArgsEnviron:     NewArgsEnviron(),
	}
	tmpVar.extSettingsArgs()
	return tmpVar
}

func (v *Variables) extSettingsArgs() {
	// Look for priority for ServerAddress
	tmpAddress := v.ArgsEnviron.GetSrvAddr()
	if tmpAddress != "" {
		v.ServerAddress = tmpAddress
	} else {
		v.ServerAddress = v.ArgsCommandLine.GetSrvAddr()
	}

	// Look for priority for BaseURL
	tmpURL := v.ArgsEnviron.GetBaseURL()
	if tmpURL != "" {
		v.BaseURL = tmpURL
	} else {
		v.BaseURL = v.ArgsCommandLine.GetBaseURL()
	}

	// Look for priority for PathToStore
	tmpPath := v.ArgsEnviron.GetPathToStore()
	if tmpPath != "" {
		v.PathToStore = tmpPath
	} else {
		v.PathToStore = v.ArgsCommandLine.GetPathToStore()
	}

	// Look for priority for DB
	tmpDB := v.ArgsEnviron.GetDB()
	if tmpDB != "" {
		v.DB = tmpDB
	} else {
		v.DB = v.ArgsCommandLine.GetDB()
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

func (v *Variables) GetLogFile() (LogFile string) {
	return v.ArgsEnviron.LogFile
}

func (v *Variables) GetDB() (DB string) {
	return v.DB
}

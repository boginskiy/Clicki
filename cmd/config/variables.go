package config

type Variables struct {
	ServerAddress   string
	BaseURL         string
	ArgsCommandLine *ArgsCommandLine
	ArgsEnviron     *ArgsEnviron
}

func NewVariables() *Variables {
	tmpVar := &Variables{
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
}

func (v *Variables) GetSrvAddr() (ServerAddress string) {
	return v.ServerAddress
}

func (v *Variables) GetBaseURL() (BaseURL string) {
	return v.BaseURL
}

func (v *Variables) GetNameLogInfo() (NameLogInfo string) {
	return v.ArgsEnviron.NameLogInfo
}

func (v *Variables) GetNameLogFatal() (NameLogFatal string) {
	return v.ArgsEnviron.NameLogFatal
}

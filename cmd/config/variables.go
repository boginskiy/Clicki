package config

type Variables struct {
	ServerAddress string
	BaseURL       string
}

func NewVariables() *Variables {
	tmpVar := &Variables{}
	tmpVar.settingsArgs()
	return tmpVar
}

func (v *Variables) checkCondition(params1, params2 Variabler) {
	// Look for priority for ServerAddress
	tmpAddress := params2.GetSrvAddr()
	if tmpAddress != "" {
		v.ServerAddress = tmpAddress
	} else {
		v.ServerAddress = params1.GetSrvAddr()
	}

	// Look for priority for BaseURL
	tmpUrl := params2.GetBaseUrl()
	if tmpUrl != "" {
		v.BaseURL = tmpUrl
	} else {
		v.BaseURL = params1.GetBaseUrl()
	}
}

func (v *Variables) settingsArgs() {
	// Create
	argscli := NewArgsCommandLine()
	argsenv := NewArgsEnviron()
	// Check
	v.checkCondition(argscli, argsenv)
}

func (v *Variables) GetSrvAddr() (ServerAddress string) {
	return v.ServerAddress
}

func (v *Variables) GetBaseUrl() (BaseURL string) {
	return v.BaseURL
}

package config

type Variables struct {
	Server_address string
	Base_url       string
}

func NewVariables() *Variables {
	tmpVar := &Variables{}
	tmpVar.settingsArgs()
	return tmpVar
}

func (v *Variables) checkCondition(params1, params2 Variabler) {
	// Look for priority for Server_address
	tmpAddress := params2.GetSrvAddr()
	if tmpAddress != "" {
		v.Server_address = tmpAddress
	} else {
		v.Server_address = params1.GetSrvAddr()
	}

	// Look for priority for Base_url
	tmpUrl := params2.GetBaseUrl()
	if tmpUrl != "" {
		v.Base_url = tmpUrl
	} else {
		v.Base_url = params1.GetBaseUrl()
	}
}

func (v *Variables) settingsArgs() {
	// Create
	argscli := NewArgsCommandLine()
	argsenv := NewArgsEnviron()
	// Check
	v.checkCondition(argscli, argsenv)
}

func (v *Variables) GetSrvAddr() (Server_address string) {
	return v.Server_address
}

func (v *Variables) GetBaseUrl() (Base_url string) {
	return v.Base_url
}

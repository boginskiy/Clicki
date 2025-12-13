package config

func (ac *ArgsCLI) Reset() {
	if ac == nil {
		return
	}

	ac.ServerAddress = ""

	ac.PathToStore = ""

	ac.AuditFile = ""

	ac.AuditURL = ""

	ac.BaseURL = ""

	ac.DB = ""

}

func (ae *ArgsENV) Reset() {
	if ae == nil {
		return
	}

	ae.ServerAddress = ""

	ae.PathToStore = ""

	ae.DB = ""

	ae.BaseURL = ""

	ae.AuditFile = ""

	ae.AuditURL = ""

	ae.LogFile = ""

	ae.MaxRetries = 0

	ae.TokenLiveTime = 0

	ae.CokiLiveTime = 0

	ae.NameCoki = ""

	ae.SecretKey = ""

	ae.SoftDeleteTime = 0

	ae.HardDeleteTime = 0

}

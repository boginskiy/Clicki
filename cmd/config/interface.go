package config

type Argsmenter interface {
	ParseFlags()
}

type ArgsCLIGetter interface {
	GetPathToStore() (PathToStore string)
	GetSrvAddr() (ServerAddress string)
	GetBaseURL() (BaseURL string)
	GetDB() (DB string)
}

type ArgsENVGetter interface {
	GetSoftDeleteTime() (SoftDeleteTime int)
	GetHardDeleteTime() (HardDeleteTime int)
	GetTokenLiveTime() (TokenLiveTime int)
	GetCokiLiveTime() (CokiLiveTime int)
	GetSecretKey() (SecretKey string)
	GetMaxRetries() (MaxRetries int)
	GetNameCoki() (NameCoki string)
	GetLogFile() (LogFile string)
}

type VarGetter interface {
	GetSoftDeleteTime() (SoftDeleteTime int)
	GetHardDeleteTime() (HardDeleteTime int)
	GetTokenLiveTime() (TokenLiveTime int)
	GetPathToStore() (PathToStore string)
	GetCokiLiveTime() (CokiLiveTime int)
	GetSrvAddr() (ServerAddress string)
	GetSecretKey() (SecretKey string)
	GetMaxRetries() (MaxRetries int)
	GetNameCoki() (NameCoki string)
	GetLogFile() (LogFile string)
	GetBaseURL() (BaseURL string)
	GetDB() (DB string)
}

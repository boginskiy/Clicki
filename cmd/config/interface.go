package config

type Argsmenter interface {
	ParseFlags()
}

type IArgsCLI interface {
	GetPathToStore() (PathToStore string)
	GetAuditFile() (AuditFile string)
	GetAuditURL() (AuditURL string)
	GetSrvAddr() (ServerAddress string)
	GetBaseURL() (BaseURL string)
	GetDB() (DB string)
}

type IArgsENV interface {
	GetSoftDeleteTime() (SoftDeleteTime int)
	GetHardDeleteTime() (HardDeleteTime int)
	GetTokenLiveTime() (TokenLiveTime int)
	GetCokiLiveTime() (CokiLiveTime int)
	GetSecretKey() (SecretKey string)
	GetMaxRetries() (MaxRetries int)
	GetNameCoki() (NameCoki string)
	GetLogFile() (LogFile string)
}

type Config interface {
	IArgsENV
	IArgsCLI
}

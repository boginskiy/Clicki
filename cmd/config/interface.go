package config

type Argsmenter interface {
	ParseFlags()
}

// IArgsCLI - interface about args comm line interface.
type IArgsCLI interface {
	GetPathToStore() (PathToStore string)
	GetAuditFile() (AuditFile string)
	GetAuditURL() (AuditURL string)
	GetSrvAddr() (ServerAddress string)
	GetBaseURL() (BaseURL string)
	GetDB() (DB string)
}

// IArgsENV - interface about args environment.
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

// Config - interface.
type Config interface {
	IArgsENV
	IArgsCLI
}

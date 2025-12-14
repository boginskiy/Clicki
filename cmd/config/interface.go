package config

type Argsmenter interface {
	ParseFlags()
}

// ConfigCLI - interface about args comm line interface.
type ConfigCLI interface {
	GetEnableHTTPS() (EnableHTTPS string)
	GetPathToStore() (PathToStore string)
	GetAuditFile() (AuditFile string)
	GetAuditURL() (AuditURL string)
	GetSrvAddr() (ServerAddress string)
	GetBaseURL() (BaseURL string)
	GetDB() (DB string)
}

// ConfigENV - interface about args environment.
type ConfigENV interface {
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
	ConfigENV
	ConfigCLI
}

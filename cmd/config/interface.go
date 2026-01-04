package config

type Argsmenter interface {
	ParseFlags()
}

// ConfigPrioryty - interface about args comm line interface.
type ConfigPrioryty interface {
	GetTrustedSubnet() (TrustedSubnet string)
	GetEnableHTTPS() (EnableHTTPS string)
	GetPathToStore() (PathToStore string)
	GetEnableGRps() (EnableGRps string)
	GetAuditFile() (AuditFile string)
	GetAuditURL() (AuditURL string)
	GetSrvAddr() (ServerAddress string)
	GetBaseURL() (BaseURL string)
	GetDB() (DB string)
}

type ConfigLog interface {
	GetLogFile() (LogFile string)
}

type ConfigDB interface {
	GetSoftDeleteTime() (SoftDeleteTime int)
	GetHardDeleteTime() (HardDeleteTime int)
	GetMaxRetries() (MaxRetries int)
}

type ConfigAuth interface {
	GetTokenLiveTime() (TokenLiveTime int)
	GetCokiLiveTime() (CokiLiveTime int)
	GetSecretKey() (SecretKey string)
	GetNameCoki() (NameCoki string)
}

// Config - interface.
type Config interface {
	ConfigPrioryty

	ConfigAuth
	ConfigLog
	ConfigDB
}

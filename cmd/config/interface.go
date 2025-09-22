package config

type Argsmenter interface {
	ParseFlags()
}

type ArgsGetter interface {
	GetPathToStore() (PathToStore string)
	GetSrvAddr() (ServerAddress string)
	GetBaseURL() (BaseURL string)
	GetDB() (DB string)
}

type VarGetter interface {
	GetPathToStore() (PathToStore string)
	GetSrvAddr() (ServerAddress string)
	GetMaxRetries() (MaxRetries int)
	GetLogFile() (LogFile string)
	GetBaseURL() (BaseURL string)
	GetDB() (DB string)
}

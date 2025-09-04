package config

type Argsmenter interface {
	ParseFlags()
}

type ArgsGetter interface {
	GetPathToStore() (PathToStore string)
	GetSrvAddr() (ServerAddress string)
	GetBaseURL() (BaseURL string)
}

type VarGetter interface {
	GetNameLogFatal() (NameLogFatal string)
	GetPathToStore() (PathToStore string)
	GetNameLogInfo() (NameLogInfo string)
	GetSrvAddr() (ServerAddress string)
	GetBaseURL() (BaseURL string)
}

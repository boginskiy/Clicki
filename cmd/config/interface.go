package config

type Argsmenter interface {
	ParseFlags()
}

type ArgsGetter interface {
	GetSrvAddr() (ServerAddress string)
	GetBaseURL() (BaseURL string)
}

type VarGetter interface {
	GetNameLogFatal() (NameLogFatal string)
	GetNameLogInfo() (NameLogInfo string)
	GetSrvAddr() (ServerAddress string)
	GetBaseURL() (BaseURL string)
}

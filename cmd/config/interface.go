package config

type Argsmenter interface {
	ParseFlags()
}

type Variabler interface {
	GetSrvAddr() (ServerAddress string)
	GetBaseURL() (BaseURL string)
}

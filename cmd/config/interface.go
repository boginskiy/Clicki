package config

type Argsmenter interface {
	ParseFlags()
}

type Variabler interface {
	GetSrvAddr() (ServerAddress string)
	GetBaseUrl() (BaseURL string)
}

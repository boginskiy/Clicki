package config

type Argsmenter interface {
	ParseFlags()
}

type Variabler interface {
	GetSrvAddr() (Server_address string)
	GetBaseUrl() (Base_url string)
}

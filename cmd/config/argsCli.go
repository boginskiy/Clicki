package config

import "flag"

type ArgsCommandLine struct {
	ServerAddress string // StartPort is the port for start application
	BaseURL       string // ResultPort is the port after changing
}

func NewArgsCommandLine() *ArgsCommandLine {
	ArgsCLI := new(ArgsCommandLine)
	ArgsCLI.ParseFlags()
	return ArgsCLI
}

func (c *ArgsCommandLine) ParseFlags() {
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "Result adress for application")
	flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "Start adress for application")
	flag.Parse()
}

func (c *ArgsCommandLine) GetSrvAddr() (ServerAddress string) {
	return c.ServerAddress
}

func (c *ArgsCommandLine) GetBaseUrl() (BaseURL string) {
	return c.BaseURL
}

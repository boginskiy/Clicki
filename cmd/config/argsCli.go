package config

import "flag"

type ArgsCommandLine struct {
	Server_address string // StartPort is the port for start application
	Base_url       string // ResultPort is the port after changing
}

func NewArgsCommandLine() *ArgsCommandLine {
	ArgsCLI := new(ArgsCommandLine)
	ArgsCLI.ParseFlags()
	return ArgsCLI
}

func (c *ArgsCommandLine) ParseFlags() {
	flag.StringVar(&c.Base_url, "b", "http://localhost:8080", "Result adress for application")
	flag.StringVar(&c.Server_address, "a", "localhost:8080", "Start adress for application")
	flag.Parse()
}

func (c *ArgsCommandLine) GetSrvAddr() (Server_address string) {
	return c.Server_address
}

func (c *ArgsCommandLine) GetBaseUrl() (Base_url string) {
	return c.Base_url
}

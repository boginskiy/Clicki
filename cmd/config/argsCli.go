package config

import "flag"

type ArgsCommandLine struct {
	ServerAddress string // StartPort is the port for start application
	PathToStore   string // PathToStore is the path to store URL
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
	flag.StringVar(&c.PathToStore, "f", "store", "Path to file of store URL")
	flag.Parse()
}

func (c *ArgsCommandLine) GetSrvAddr() (ServerAddress string) {
	return c.ServerAddress
}

func (c *ArgsCommandLine) GetBaseURL() (BaseURL string) {
	return c.BaseURL
}

func (c *ArgsCommandLine) GetPathToStore() (PathToStore string) {
	return c.PathToStore
}

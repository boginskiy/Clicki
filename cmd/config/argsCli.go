package config

import "flag"

type ArgsCLI struct {
	ServerAddress string // StartPort is the port for start application
	PathToStore   string // PathToStore is the path to store URL
	BaseURL       string // ResultPort is the port after changing
}

func NewArgsCLI() *ArgsCLI {
	args := new(ArgsCLI)
	args.ParseFlags()
	return args
}

func (c *ArgsCLI) ParseFlags() {
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "Result adress for application")
	flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "Start adress for application")
	flag.StringVar(&c.PathToStore, "f", "store", "Path to file of store URL")
	flag.Parse()
}

func (c *ArgsCLI) GetSrvAddr() (ServerAddress string) {
	return c.ServerAddress
}

func (c *ArgsCLI) GetBaseURL() (BaseURL string) {
	return c.BaseURL
}

func (c *ArgsCLI) GetPathToStore() (PathToStore string) {
	return c.PathToStore
}

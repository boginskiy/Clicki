package config

import "flag"

type Argsmenter interface {
	ParseFlags()
}

type ArgumentsCLI struct {
	StartPort  string // StartPort is the port for start application
	ResultPort string // ResultPort is the port after changing
}

func NewArgumentsCLI() *ArgumentsCLI {
	ArgsCLI := new(ArgumentsCLI)
	ArgsCLI.ParseFlags()
	return ArgsCLI
}

func (a *ArgumentsCLI) ParseFlags() {
	flag.StringVar(&a.ResultPort, "b", "http://localhost:8080", "Result adress for application")
	flag.StringVar(&a.StartPort, "a", "localhost:8080", "Start adress for application")
	flag.Parse()
}

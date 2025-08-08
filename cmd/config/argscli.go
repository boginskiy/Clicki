package config

import "flag"

type Argsmenter interface {
	ParseFlags()
}

type ArgumentsCLI struct {
	StartPort  string // Порт запуска приложения
	ResultPort string // Порт на который перенаправим новый URL
}

func ParseFlags() *ArgumentsCLI {
	ArgsCLI := new(ArgumentsCLI)
	flag.StringVar(&ArgsCLI.ResultPort, "b", "http://localhost:8080", "Result adress for application")
	flag.StringVar(&ArgsCLI.StartPort, "a", "localhost:8080", "Start adress for application")
	flag.Parse()
	return ArgsCLI
}

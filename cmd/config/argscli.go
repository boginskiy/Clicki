package config

import "flag"

var ArgsCLI *ArgumentsCLI

type ArgumentsCLI struct {
	StartPort  string // Порт запуска приложения
	ResultPort string // Порт на который перенаправим новый URL
}

func ParseFlags() {
	ArgsCLI = new(ArgumentsCLI)
	flag.StringVar(&ArgsCLI.StartPort, "a", "localhost:8080", "Start adress for application")
	flag.StringVar(&ArgsCLI.ResultPort, "b", "http://localhost:8080", "Result adress for application")
	flag.Parse()
}

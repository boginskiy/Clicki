package config

import "flag"

var defaultPort = "8080"
var ArgsCLI *ArgumentsCLI

type ArgumentsCLI struct {
	StartPort  string
	ResultPort string
}

func ParseFlags() {
	ArgsCLI = new(ArgumentsCLI)
	flag.StringVar(&ArgsCLI.StartPort, "a", "", "Start port for application")
	flag.StringVar(&ArgsCLI.ResultPort, "b", "", "Result port for application")
	flag.Parse()

	// Обработка особых случаев

	switch {
	case ArgsCLI.StartPort != "" && ArgsCLI.ResultPort == "":
		// флага '-b' нет, уcтанавливаем значение как у флага '-a'
		ArgsCLI.ResultPort = ArgsCLI.StartPort
	case ArgsCLI.StartPort == "" && ArgsCLI.ResultPort != "":
		// флага '-a' нет, но есть флаг '-b'. Флаг '-a' задается с defaultPort
		ArgsCLI.StartPort = defaultPort
	case ArgsCLI.StartPort == "" && ArgsCLI.ResultPort == "":
		ArgsCLI.StartPort = defaultPort
		ArgsCLI.ResultPort = defaultPort
	}
}

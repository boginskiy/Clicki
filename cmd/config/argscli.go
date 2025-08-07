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
	flag.StringVar(&ArgsCLI.ResultPort, "b", "http://localhost:8080/", "Result adress for application")
	flag.Parse()

	// // Обработка особых случаев
	// switch {

	// // флага '-b' нет, уcтанавливаем значение как у флага '-a'
	// case ArgsCLI.StartPort != "" && ArgsCLI.ResultPort == "":
	// 	ArgsCLI.ResultPort = ArgsCLI.StartPort

	// // флага '-a' нет, но есть флаг '-b'. Флаг '-a' задается с defaultPort
	// case ArgsCLI.StartPort == "" && ArgsCLI.ResultPort != "":
	// 	ArgsCLI.StartPort = defaultPort
	// 	ArgsCLI.IsCh = true

	// case ArgsCLI.StartPort == "" && ArgsCLI.ResultPort == "":
	// 	ArgsCLI.StartPort = defaultPort
	// 	ArgsCLI.ResultPort = defaultPort

	// default:
	// 	ArgsCLI.IsCh = true
	// }
}

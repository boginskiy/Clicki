package main

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/server"
)

func main() {
	// Обработчик командной строки
	config.ParseFlags()

	if err := server.Run(); err != nil {
		panic(err)
	}
}

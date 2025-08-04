package main

import (
	"github.com/boginskiy/Clicki/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		panic(err)
	}
}

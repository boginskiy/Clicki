package main

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/server"
)

func main() {
	Variables := config.NewVariables()
	server.Run(Variables)
}

package main

import (
	"log"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/server"
)

func main() {
	Variables := config.NewVariables()

	if err := server.Run(Variables); err != nil {
		log.Fatal(err)
	}
}

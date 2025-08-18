package main

import (
	"log"

	"github.com/boginskiy/Clicki/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

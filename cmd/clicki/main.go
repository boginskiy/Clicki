package main

import (
	"github.com/boginskiy/Clicki/internal/router"
)

func main() {
	if err := router.Run(); err != nil {
		panic(err)
	}
}

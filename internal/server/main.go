package server

import (
	"fmt"
	"net/http"

	c "github.com/boginskiy/Clicki/cmd/config"
	db "github.com/boginskiy/Clicki/internal/db"
	r "github.com/boginskiy/Clicki/internal/router"
	s "github.com/boginskiy/Clicki/internal/service"
)

func Run() error {
	// agrs - атрибуты командной строки
	argsCLI := c.ParseFlags()

	// db - слой базы данных 'DbStore'
	db := db.NewDbStore()

	// shortingURL - слой с бизнес логикой сервиса 'ShorteningURL'
	shortingURL := s.NewShorteningURL(db)

	fmt.Printf("The server has started on port %s\n", argsCLI.StartPort)
	return http.ListenAndServe(argsCLI.StartPort, r.Router(shortingURL, argsCLI))
}

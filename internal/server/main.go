package server

import (
	"fmt"
	"net/http"

	c "github.com/boginskiy/Clicki/cmd/config"
	db "github.com/boginskiy/Clicki/internal/db"
	p "github.com/boginskiy/Clicki/internal/preparation"
	r "github.com/boginskiy/Clicki/internal/router"
	s "github.com/boginskiy/Clicki/internal/service"
	v "github.com/boginskiy/Clicki/internal/validation"
)

func Run() error {
	// agrs - атрибуты командной строки
	argsCLI := c.NewArgumentsCLI()

	extraFuncer := p.NewExtraFunc() // extraFuncer - дополнительные возможности
	checker := v.NewChecker()       // checker - валидация данных
	db := db.NewDBStore()           // db - слой базы данных 'DBStore'

	// shortingURL - слой с бизнес логикой сервиса 'ShorteningURL'
	shortingURL := s.NewShorteningURL(db, checker, extraFuncer)

	fmt.Printf("The server has started on port %s\n", argsCLI.StartPort)
	return http.ListenAndServe(argsCLI.StartPort, r.Router(shortingURL, argsCLI))
}

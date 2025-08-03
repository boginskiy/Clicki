package router

import (
	"fmt"
	"net/http"

	hdl "github.com/boginskiy/Clicki/internal/handler"
)

func serveMux() *http.ServeMux {
	tmpHandler := hdl.NewHandlerForURL()
	tmpMux := http.NewServeMux()
	tmpMux.Handle("/", tmpHandler)
	return tmpMux
}

func Run() error {
	mux := serveMux()

	defer fmt.Println("The server has stopped on port 8080")
	fmt.Println("The server has started on port 8080")
	return http.ListenAndServe(":8080", mux)
}

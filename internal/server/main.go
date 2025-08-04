package server

import (
	"fmt"
	"net/http"

	"github.com/boginskiy/Clicki/internal/router"
)

func Run() error {
	mux := router.ServeMux()

	fmt.Println("The server has started on port 8080")
	return http.ListenAndServe(":8080", mux)
}

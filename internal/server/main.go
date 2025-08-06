package server

import (
	"fmt"
	"net/http"

	r "github.com/boginskiy/Clicki/internal/router"
)

func Run() error {
	routerStart := r.Router()
	fmt.Println("The server has started on port 8080")
	return http.ListenAndServe(":8080", routerStart)
}

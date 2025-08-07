package server

import (
	"fmt"
	"net/http"

	c "github.com/boginskiy/Clicki/cmd/config"
	r "github.com/boginskiy/Clicki/internal/router"
)

func Run() error {
	routerStart := r.Router()
	fmt.Printf("The server has started on port %s", c.ArgsCLI.StartPort)
	return http.ListenAndServe(fmt.Sprintf(":%s", c.ArgsCLI.StartPort), routerStart)
}

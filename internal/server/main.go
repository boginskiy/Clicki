package server

import (
	"fmt"
	"net/http"

	c "github.com/boginskiy/Clicki/cmd/config"
	r "github.com/boginskiy/Clicki/internal/router"
)

func Run() error {
	fmt.Printf("The server has started on port %s\n", c.ArgsCLI.StartPort)
	return http.ListenAndServe(c.ArgsCLI.StartPort, r.Router())
}

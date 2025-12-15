package server

import (
	"fmt"
	"os"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
	mv "github.com/boginskiy/Clicki/internal/middleware"
	"github.com/boginskiy/Clicki/internal/router"
)

func Run(cfg config.Config, logg logg.Logger, router router.Router, mdlwere mv.Middleware) {
	if cfg.GetEnableHTTPS() != "1" {
		fmt.Fprintf(os.Stdout, "Protocol:      %s\n", "HTTP")
		NewServ(cfg, logg).Run(router, mdlwere)

	} else {
		fmt.Fprintf(os.Stdout, "Protocol:      %s\n", "HTTPS")
		NewServS(cfg, logg).Run(router, mdlwere)
	}
}

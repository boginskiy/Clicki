package server

import (
	"fmt"
	"os"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
	mv "github.com/boginskiy/Clicki/internal/middleware"
	"github.com/boginskiy/Clicki/internal/router"
	"github.com/boginskiy/Clicki/internal/rpc"
)

func RunHTTP(cfg config.Config, logg logg.Logger, router router.Router, mdlwere mv.Middleware) {
	if cfg.GetEnableHTTPS() == "1" {
		fmt.Fprintf(os.Stdout, "Protocol:      %s\n", "HTTPS")
		NewServS(cfg, logg, router.Run(mdlwere)).Run()

	} else {
		fmt.Fprintf(os.Stdout, "Protocol:      %s\n", "HTTP")
		NewServ(cfg, logg, router.Run(mdlwere)).Run()
	}
}

func RunGRPC(cfg config.Config, logg logg.Logger, service rpc.ShortenerServiceServer, intercept mv.Interceptor) {
	if cfg.GetEnableGRps() == "1" {
		fmt.Fprintf(os.Stdout, "Protocol:      %s\n", "gRPC")
		NewServgRPC(cfg, logg, service, intercept).Run()
	}
}

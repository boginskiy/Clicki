package server

import (
	"net/http"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
	mv "github.com/boginskiy/Clicki/internal/middleware"
	"github.com/boginskiy/Clicki/internal/router"
)

type Server interface {
	Run(router router.Router, middleware mv.Middleware)
}

type Serv struct {
	Cfg  config.Config
	Logg logg.Logger
}

func NewServ(config config.Config, logger logg.Logger) *Serv {
	return &Serv{
		Cfg:  config,
		Logg: logger,
	}
}

func (s *Serv) Run(router router.Router, mdlwere mv.Middleware) {
	s.Logg.RaiseFatal(
		http.ListenAndServe(s.Cfg.GetSrvAddr(), router.Run(mdlwere)),
		"server has not started", nil)
}

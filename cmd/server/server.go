package server

import (
	"net/http"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
	mv "github.com/boginskiy/Clicki/internal/middleware"
	"github.com/boginskiy/Clicki/internal/router"
	"golang.org/x/crypto/acme/autocert"
)

var DIRCACHE = "cache-dir"
var HOST = "mysite.ru"

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

func (s *Serv) StartWithAutocert(nameSite string, handler http.Handler) {
	s.Logg.RaiseFatal(http.Serve(autocert.NewListener(nameSite), handler),
		"server has not started", nil)
}

// SettingsManager конструируем менеджер TLS-сертификатов.
func (s *Serv) SettingsManager(dirCache string, hosts ...string) *autocert.Manager {
	return &autocert.Manager{
		Cache:      autocert.DirCache(dirCache),      // директория для хранения сертификатов.
		Prompt:     autocert.AcceptTOS,               // функция, принимающая Terms of Service издателя сертификатов.
		HostPolicy: autocert.HostWhitelist(hosts...), // перечень доменов, для которых будут поддерживаться сертификаты.
	}
}

// SettingsServerTLS конструируем сервер с поддержкой TLS
func (s *Serv) SettingsServerTLS(manager *autocert.Manager, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    ":443",
		Handler: handler,
		// для TLS-конфигурации используем менеджер сертификатов
		TLSConfig: manager.TLSConfig(),
	}
}

// StartWithCustomAutocert. Cервер запускать с правами администратора и на хосте, имеющем «белый» IP-адрес.
func (s *Serv) StartWithCustomAutocert(handler http.Handler) {
	autocertManager := s.SettingsManager(DIRCACHE, HOST)
	server := s.SettingsServerTLS(autocertManager, handler)

	s.Logg.RaiseFatal(
		server.ListenAndServeTLS("", ""),
		"server has not started", nil)
}

func (s *Serv) StartServeTLS(handler http.Handler) {
	s.Logg.RaiseFatal(
		http.ListenAndServeTLS(s.Cfg.GetSrvAddr(), handler),
		"server has not started", nil)
}

func (s *Serv) Run(router router.Router, mdlwere mv.Middleware) {
	if s.Cfg.GetEnableHTTPS() == "1" {
		// s.StartWithAutocert(HOST, router.Run(mdlwere))
		// s.StartWithCustomAutocert(router.Run(mdlwere))
		s.StartServeTLS(router.Run(mdlwere))

	} else {
		s.Logg.RaiseFatal(
			http.ListenAndServe(s.Cfg.GetSrvAddr(), router.Run(mdlwere)),
			"server has not started", nil)
	}
}
